package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mariajdab/pay-gateway/internal/AcqBank"
	"github.com/mariajdab/pay-gateway/internal/cards"
	"github.com/mariajdab/pay-gateway/internal/customers"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/merchant"
	"github.com/mariajdab/pay-gateway/internal/payments"
)

// statuses of transactions according "bank documentation"
var txnStatusMap = map[uint16]string{
	2020: allowed,
	2222: declined,
}

// txn statuses
const (
	allowed  = "txn-allowed"
	declined = "txn-declined"
)

type paymentResp struct {
	StatusPayment string `json:"status_payment"`
	TxnUUID       string `json:"txn_uuid"`
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
	merchantCode := ctx.Param("merchant_code")
	fmt.Println(merchantCode, "mercado")

	if merchantCode == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	payReq := new(entity.PaymentRequest)
	err := ctx.BindJSON(payReq)
	fmt.Println(payReq, "salito")
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	validMer, err := h.ucMerchant.ValidateMerchant(merchantCode)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if validMer != entity.AllowedMerchant {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var response = paymentResp{}
	paymtValidateResp, err := h.ucPayment.ValidatePaymentReq(*payReq)

	if err != nil {
		log.Printf("error in ValidatePaymentReq: %s", err)
		response.StatusPayment = "Failed"
		ctx.JSON(http.StatusOK, response)
		return
	}

	if paymtValidateResp.Status == entity.PendingBankValidation {
		// we need to validate the card with the bank
		cardData := entity.CardData{
			CardInfo:  payReq.CardInfo,
			OwnerInfo: payReq.CustomerData,
		}

		cardBankRes, err := validateCardWithBank(cardData)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if cardBankRes.Valid {
			paymtValidateResp.CardUUIDBank = cardBankRes.CardUUID

			// save card owner info
			custID, err := h.ucCustomer.CreateCustomer(payReq.CustomerData)
			if err != nil {
				log.Println("internal error during CreateCustomer", err)
			}

			// save the card info
			if err := h.ucCard.SaveCardInfo(
				paymtValidateResp.CardTk,
				cardBankRes.CardUUID,
				payReq.CardInfo.ExpDate,
				custID); err != nil {
				log.Println("internal error during SaveCardInfo", err)
			}

		} else { // the bank status of the card is invalid
			response.StatusPayment = "failed"
			ctx.JSON(http.StatusOK, response)
		}
	}

	// we have the cardUUIDBank for new or known card
	txn := entity.TxnInfo{
		CardUUID: paymtValidateResp.CardUUIDBank,
		Amount:   payReq.BillingAmount,
		Currency: payReq.Currency,
	}

	status, err := checkTxnWithBank(txn)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if status == allowed {
		// save payment in the db
		transaction := entity.Transaction{
			BillingAmount: payReq.BillingAmount,
			CreatedAt:     time.Now(),
			Status:        "success",
			CardTk:        paymtValidateResp.CardTk,
			MerchantCode:  merchantCode,
		}

		txnUUID, err := h.ucPayment.SavePaymentInfo(transaction)
		if err != nil {
			log.Println("internal error during SavePaymentInfo", err)
		}

		response.StatusPayment = "Success"
		response.TxnUUID = txnUUID
	} else {
		response.StatusPayment = "Failed"
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) RetrievePayment(ctx *gin.Context) {
	merchantCode := ctx.Param("merchant_code")
	fmt.Println(merchantCode, "mercado")

	paymentUUID := ctx.Param("payment_uuid")
	fmt.Println(merchantCode, "mercado")

	if merchantCode == "" || paymentUUID == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	paymentDetail, err := h.ucPayment.PaymentDetailByTxnUUID(paymentUUID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	} else if err == nil && paymentDetail.MerchantCode == "" { // the payment was not found
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// merchant_code param must match with merchant code of the payment in db
	if paymentDetail.MerchantCode != merchantCode {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, paymentDetail)
}

func validateCardWithBank(cardData entity.CardData) (AcqBank.CardValidResp, error) {
	url := "/cards/validate"

	body, err := json.Marshal(cardData)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Println("error during validation card request")
		return AcqBank.CardValidResp{}, err
	}

	defer res.Body.Close()

	cardValidResp := AcqBank.CardValidResp{}

	if err := json.NewDecoder(res.Body).Decode(&cardValidResp); err != nil {
		log.Println("internal error during decode card validation response", err)
		return AcqBank.CardValidResp{}, err
	}

	return cardValidResp, nil
}

func checkTxnWithBank(txnInfo entity.TxnInfo) (string, error) {
	url := "/transactions/validate"

	fmt.Println(txnInfo, "venusa")

	txn := entity.TxnInfo{
		CardUUID: txnInfo.CardUUID,
		Amount:   txnInfo.Amount,
		Currency: txnInfo.Currency,
	}

	body, err := json.Marshal(&txn)
	if err != nil {
		log.Println("internal error during marshal txnInfo", err)
		return "", err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("error making HTTP request:", err)
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	txnResp := AcqBank.TxnResp{}

	if err := json.NewDecoder(res.Body).Decode(&txnResp); err != nil {
		log.Println("internal error during decode txn response", err)
		return "", err
	}

	fmt.Println(txnResp, "lulita")

	status, ok := txnStatusMap[txnResp.Status]
	if !ok {
		log.Println("internal error tx code from Bank is unknown")
		return "", errors.New("unknown bank code")
	}

	return status, nil
}
