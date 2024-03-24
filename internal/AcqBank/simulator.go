package AcqBank

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"net/http"
	"strconv"
)

const (
	approved = 2020
	declined = 9999
	pending  = 8585
)

type CheckTxResp struct {
	Status uint16 `json:"status_txn"`
}

type CardValidResp struct {
	Valid    bool   `json:"valid"`
	CardUUID string `json:"card_uuid"`
}

type Handler struct{}

func RegisterHTTPEndpoints(g *gin.Engine, handler *Handler) {

	payment := g.Group("/bank-sim.com")
	{
		payment.GET("/cards/:card_uuid", handler.checkTxn)
		payment.GET("/cards/", handler.validateCard)

	}
}

func (h *Handler) checkTxn(ctx *gin.Context) {
	cardUUID := ctx.Param("card_uuid")
	if cardUUID == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	status := getTxnStatus(cardUUID)

	statusTxn := CheckTxResp{status}

	ctx.JSON(http.StatusOK, statusTxn)

}

// validateCard handle the simulation logic to check if a card
// is a valid bank card and return the cardUUID for better communication
// with the payment gateway
func (h *Handler) validateCard(ctx *gin.Context) {
	card := new(entity.Card)
	if err := ctx.BindJSON(card); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
	}

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

		ctx.JSON(http.StatusOK, cardResp)
	} else {
		cardResp.Valid = false
		ctx.JSON(http.StatusOK, cardResp)
	}
}

func getTxnStatus(cardUUID string) uint16 {
	lastCharCard := string(cardUUID[len(cardUUID)-1])

	if isLastCharacterInt(lastCharCard) {
		if isOddNumber(lastCharCard) {
			return pending
		} else {
			return declined
		}
	}
	// if the last character in the uuid is a letter, then approve the txn
	return approved
}

func isLastCharacterInt(lastChar string) bool {
	// if there is an error it should be character and not a number
	if _, err := strconv.Atoi(lastChar); err == nil {
		return true
	} else {
		return false
	}
}

func isOddNumber(lastChar string) bool {
	// should not be an error
	v, _ := strconv.Atoi(lastChar)
	if v%2 > 0 {
		return true
	}
	return false
}
