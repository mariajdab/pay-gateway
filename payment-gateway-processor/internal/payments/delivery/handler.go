package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mariajdab/pay-gateway/internal/cards"
	"github.com/mariajdab/pay-gateway/internal/customers"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/entity/util"
	"github.com/mariajdab/pay-gateway/internal/merchant"
	"github.com/mariajdab/pay-gateway/internal/payments"
)

// txn statuses
const (
	// statuses of transactions according "bank documentation"

	allowedTxByBank  = 2020
	declinedTxByBank = 2222

	declinedRefundByBank = 4444
	acceptedRefundByBank = 2030
)

type paymentResp struct {
	StatusPayment string `json:"status_payment"`
	TxnUUID       string `json:"txn_uuid"`
	Reason        string `json:"reason,omitempty"`
}

type refundResp struct {
	StatusRefund string `json:"status_refund"`
	Reason       string `json:"reason,omitempty"`
}

type RefundReq struct {
	TxnUUID      string  `json:"txn_uuid"`
	Amount       float32 `json:"amount"`
	MerchantCode string  `json:"merchant_code"`
}

type Handler struct {
	ucPayment  payments.UseCasePayment
	ucCard     cards.UseCaseCard
	ucMerchant merchant.UseCaseMerchant
	ucCustomer customers.UseCaseCustomer
}

func NewHandler(ucp payments.UseCasePayment, ucc cards.UseCaseCard, ucm merchant.UseCaseMerchant, uccust customers.UseCaseCustomer) *Handler {
	return &Handler{
		ucPayment:  ucp,
		ucCard:     ucc,
		ucMerchant: ucm,
		ucCustomer: uccust,
	}
}

func (h *Handler) ProcessPayment(ctx *gin.Context) {
	payReq := new(entity.PaymentRequest)
	err := ctx.BindJSON(payReq)

	if err != nil {
		log.Println("invalid request sent to ProcessPayment")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var response = paymentResp{}
	paymtValidateResp, err := h.ucPayment.ValidatePaymentReq(*payReq)

	if err != nil {
		log.Println("payment validation request failed, error", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cardData := entity.CardData{
		CardInfo:  payReq.CardInfo,
		OwnerInfo: payReq.CustomerData,
	}

	// simulate encryption data card to be sent to the bank
	cardEncrypted, err := util.EncryptCardData(cardData)
	if err != nil {
		log.Println("internal error during encrypt card data", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if paymtValidateResp.Status == entity.PendingBankValidation { // we need to validate the card with the bank
		cardBankResp, err := validateCardWithBank(cardEncrypted)
		if err != nil {
			log.Println("internal error during validateCardWithBank", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if cardBankResp.Valid {
			// save card owner info
			custID, err := h.ucCustomer.CreateCustomer(payReq.CustomerData)
			if err != nil {
				// don't need to stop the server
				log.Println("internal error during CreateCustomer,", err)
			}

			// save the card info, but not sensitive information about the card
			if err := h.ucCard.SaveCardInfo(
				paymtValidateResp.CardTk, // "tokenization" of the card by the pay-processor
				payReq.CardInfo.ExpDate,
				custID); err != nil {
				// don't need to stop the server
				log.Println("internal error during SaveCardInfo", err)
			}

		} else { // the bank status of the card is invalid
			response.StatusPayment = "Error"
			response.Reason = "invalid card"
			ctx.JSON(http.StatusOK, response)
			return
		}
	}

	// we now know that the card is valid, so with can create a txn with the bank for that card
	txn := entity.TxnInfo{
		MerchantAccount:   paymtValidateResp.MerchantAccount,
		Amount:            payReq.BillingAmount,
		Currency:          payReq.Currency,
		CardInfoEncrypted: cardEncrypted,
	}

	status, err := createTxnWithBank(txn)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if status == allowedTxByBank { // only saves success payments
		// create the uuid associated with the payment that could be used by the merchant
		txnUUID := uuid.New()
		transaction := entity.Transaction{
			TxnUUID:       txnUUID.String(),
			BillingAmount: payReq.BillingAmount,
			CreatedAt:     time.Now().UTC(),
			Status:        "Success",
			CardTk:        paymtValidateResp.CardTk,
			Currency:      payReq.Currency,
			MerchantCode:  payReq.MerchantCode,
		}

		// saving not sensitive information about the card
		if err := h.ucPayment.SavePaymentInfo(transaction); err != nil {
			// don't need to stop the server
			log.Println("internal error during SavePaymentInfo", err)
		}

		response.StatusPayment = "Success"
		response.TxnUUID = txnUUID.String()
	} else {
		response.StatusPayment = "Failed"
		response.Reason = "Declined by Bank"
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) RetrievePayment(ctx *gin.Context) {
	merchantCode := ctx.Param("merchant_code")
	txnUUID := ctx.Param("payment_uuid")

	if merchantCode == "" || txnUUID == "" {
		log.Println("should be provided the merchant code and txnUUID")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	paymentDetail, err := h.ucPayment.PaymentDetail(txnUUID, merchantCode)
	if err != nil {
		log.Println("internal error during retrieve payment detail", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	} else if err == nil && paymentDetail.Status == "" { // the payment was not found
		log.Println("the req dont match any payment detail, mercahnt code and txnUUID could be not valid")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, paymentDetail)
}

func (h *Handler) RefundPayment(ctx *gin.Context) {
	refundReq := new(RefundReq)

	if err := ctx.BindJSON(refundReq); err != nil {
		log.Println("invalid request sent to RefundPayment")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	refundValidateResp, err := h.ucPayment.ValidatePaymentRefundReq(refundReq.TxnUUID, refundReq.MerchantCode)
	if err != nil {
		log.Println("internal error during ValidatePaymentRefundReq", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var response refundResp

	if !refundValidateResp.Valid {
		if refundValidateResp.Reason == entity.Refunded {
			response.StatusRefund = "Rejected"
			response.Reason = entity.Refunded
			ctx.JSON(http.StatusOK, response)
			return

		} else {
			// payment do not exists
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	// simulate call API to the bank to verify if the merchant
	// account have the amount to debit it
	refundStatus, err := createRefundWithBank(refundValidateResp.TxnInfo)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if refundStatus != acceptedRefundByBank {
		response.StatusRefund = "Rejected"
		response.Reason = entity.DeclinedRefund
		ctx.JSON(http.StatusOK, response)
	}

	if err := h.ucPayment.UpdatePaymentTxnStatus(refundReq.TxnUUID, "Refunded"); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response.StatusRefund = "Accepted"
	ctx.JSON(http.StatusOK, response)

}

func validateCardWithBank(cardDataEncrypted []byte) (entity.CardValidResp, error) {
	url := "http://bank-sim:9090/bank-sim/cards/validate"

	card := struct {
		CardInfo      []byte `json:"card_info"`
		GatwyPamtUUID string `json:"gatwy_pamt_uuid"`
	}{
		CardInfo:      cardDataEncrypted,
		GatwyPamtUUID: uuid.New().String(),
	}

	body, err := json.Marshal(&card)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("error during validation card request", err)
		return entity.CardValidResp{}, err
	}

	defer res.Body.Close()

	cardValidResp := entity.CardValidResp{}

	if err := json.NewDecoder(res.Body).Decode(&cardValidResp); err != nil {
		log.Println("internal error during decode card validation response", err)
		return entity.CardValidResp{}, err
	}

	return cardValidResp, nil
}

func createTxnWithBank(txnInfo entity.TxnInfo) (int, error) {
	url := "http://bank-sim:9090/bank-sim/transactions"

	body, err := json.Marshal(&txnInfo)
	if err != nil {
		log.Println("internal error during marshal txnInfo", err)
		return 0, err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("error making HTTP request:", err)
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	txnResp := entity.TxnResp{}

	if err := json.NewDecoder(res.Body).Decode(&txnResp); err != nil {
		log.Println("internal error during decode txn response", err)
		return 0, err
	}

	return txnResp.Status, nil
}

func createRefundWithBank(txnInfo entity.TxnInfo) (uint16, error) {
	url := "http://bank-sim:9090/bank-sim/transactions/refund"

	body, err := json.Marshal(&txnInfo)
	if err != nil {
		log.Println("internal error during marshal txnInfo", err)
		return 0, err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("error making HTTP request:", err)
		return 0, err
	}
	defer res.Body.Close()

	responseBank := entity.RefundStatus{}

	if err := json.NewDecoder(res.Body).Decode(&responseBank); err != nil {
		log.Println("internal error during decode txn response", err)
		return 0, err
	}

	return responseBank.Status, nil
}
