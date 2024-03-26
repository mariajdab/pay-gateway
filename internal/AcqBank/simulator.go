package AcqBank

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mariajdab/pay-gateway/internal/entity"
)

const (
	approved = 2020
	declined = 9999
)

type TxnResp struct {
	Status uint16 `json:"status_txn"`
}

type CardValidResp struct {
	Valid    bool   `json:"valid"`
	CardUUID string `json:"card_uuid"`
}

type Handler struct{}

func RegisterHTTPEndpoints(g *gin.Engine, handler *Handler) {
	payment := g.Group("/bank-sim")
	{
		payment.POST("/transactions/validate", handler.checkTxn)
		payment.POST("/cards/validate", handler.validateCard)
	}
}

// checkTxn handle the simulation logic to check if a transaction
// is approved, declined or pending by the bank "using" the card info
// associated with an account
func (h *Handler) checkTxn(ctx *gin.Context) {
	// TxnInfo have the necessary information to
	// check if the account associated with the card
	// has enough balance to approve or not the txn
	txnInfo := new(entity.TxnInfo)
	if err := ctx.BindJSON(txnInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}
	// amount info is not used in the simulation, just cardUUID
	cardUUID := txnInfo.CardUUID

	checkTxRp := TxnResp{}
	lastCharCard := string(cardUUID[len(cardUUID)-1])

	if isLastCharacterInt(lastCharCard) {
		checkTxRp.Status = declined // decline for even numbers
	} else { // if the last character in the cardUUID is a letter, then approve the txn
		checkTxRp.Status = approved
	}

	ctx.JSON(http.StatusOK, checkTxRp)

}

// validateCard handle the simulation logic to check if a card
// is a valid bank card, and return the cardUUID for better communication
// with the payment gateway
func (h *Handler) validateCard(ctx *gin.Context) {
	// a bank need the card info but also the owner info
	cardData := new(entity.CardData)
	if err := ctx.BindJSON(cardData); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

	card := cardData.CardInfo
	// get the last number of the card
	lastNumberStr := string(card.Number[len(card.Number)-1])
	lastNumb, err := strconv.Atoi(lastNumberStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	cardResp := CardValidResp{}
	// if last number >= 5, the card is a customer bank card valid, then
	// response OK with the cardUUID for later communication between
	// the payment gateway and the bank
	if lastNumb >= 5 {
		cardResp.CardUUID = uuid.New().String()
		cardResp.Valid = true
	} else {
		cardResp.Valid = false
	}
	ctx.JSON(http.StatusOK, cardResp)
}

func isLastCharacterInt(lastChar string) bool {
	// if there is an error it should be character and not a number
	if _, err := strconv.Atoi(lastChar); err == nil {
		return true
	} else {
		return false
	}
}
