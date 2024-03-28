package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	// tx responses
	approved = 2020
	declined = 9999
)

type TxnResp struct {
	Status int `json:"status_txn"`
}

type CardValidResp struct {
	Valid bool `json:"valid"`
}

type RefundStatus struct {
	Status uint16 `json:"refund_status"`
}
type CardValidateReq struct {
	CardInfo      []byte `json:"card_info"`
	GatwyPamtUUID string `json:"gatwy_pamt_uuid"`
}

type Handler struct{}

func RegisterHTTPEndpoints(g *gin.Engine) {
	h := NewHandler()

	payment := g.Group("/bank-sim")
	{
		payment.POST("/transactions", h.createTxn)             // for create a new transaction with the card info
		payment.POST("/cards/validate", h.validateCard)        // for validate the card info in the bank
		payment.POST("transactions/refund", h.createTxnRefund) // for create a new refund transaction
	}
}

// createTxn handle the simulation logic to create a transaction and return if it's
// approved, declined or pending by the bank "using" the card info
// associated with an account
func (h *Handler) createTxn(ctx *gin.Context) {
	// TxnInfo have the necessary information to
	// check if the account associated with the card
	// has enough balance to approve or not the tx.
	txnInfo := new(TxnInfo)
	if err := ctx.BindJSON(txnInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

	log.Printf("==== Starting Bank simulator: createTx in bank with TxnInfo ====: %+v", txnInfo)

	// amount info is used in the simulation to check is it is even or odd and then assign the txn status
	amount := txnInfo.Amount

	checkTxRp := TxnResp{}

	if int(amount)%2 > 0 {
		log.Println("===== Bank Simulator should be declined due last char in cardUUI token is int =====")
		checkTxRp.Status = declined // decline for odd amount
	} else { //  approve the txn for even amount
		checkTxRp.Status = approved
	}

	log.Printf("===== Bank Simulator: create response for createTxn: ===== %+v", checkTxRp)
	ctx.JSON(http.StatusOK, checkTxRp)

}

// validateCard handle the simulation logic to check if a card
// is a valid bank card
func (h *Handler) validateCard(ctx *gin.Context) {
	cardData := new(CardValidateReq)
	if err := ctx.BindJSON(cardData); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}
	log.Printf("==== Starting Bank simulator: validateCard in bank with TxnInfo ====: %+v", cardData)

	cardResp := CardValidResp{}
	pymtGtwUUID := cardData.GatwyPamtUUID
	lastCharPymGtwUUID := string(pymtGtwUUID[len(pymtGtwUUID)-1])

	if isLastCharacterInt(lastCharPymGtwUUID) {
		cardResp.Valid = false // decline for numbers
	} else { // if the last character in the pymtGtwUUID is a letter, then approve the txn
		cardResp.Valid = true
	}

	log.Printf("======== Bank Simulator: card validation response ======== %+v", cardResp)
	ctx.JSON(http.StatusOK, cardResp)
}

// createTxnRefund simulate the refund response from the bank, in this case
// the code just check is the amount is even or nor to return the status of the refund
func (h *Handler) createTxnRefund(ctx *gin.Context) {
	txnInfo := new(TxnInfo)
	if err := ctx.BindJSON(txnInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

	var refundResp RefundStatus

	if int(txnInfo.Amount)%2 > 0 {
		refundResp.Status = 4444 // denied
	} else {
		refundResp.Status = 2030 // accepted
	}
	ctx.JSON(http.StatusOK, refundResp)
}

func isLastCharacterInt(lastChar string) bool {
	// if there is an error it should be character and not a number
	if _, err := strconv.Atoi(lastChar); err == nil {
		return true
	} else {
		return false
	}
}
