package AcqBank

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mariajdab/pay-gateway/internal/entity"
)

const (
	// txn statuses
	allowed  = "txn-allowed"
	declined = "txn-declined"
	pending  = "txn-pending"
)

var txnStatusMap = map[int]string{
	0: allowed,
	1: pending,
	2: declined,
}

type Handler struct{}

func RegisterHTTPEndpoints(g *gin.Engine, handler *Handler) {

	payment := g.Group("/bank-sim")
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

	rand.Seed(time.Now().UnixNano())
	randomTxStatus := rand.Intn(3)

	txnStatus := struct {
		StatusTxn string `json:"status_txn"`
	}{StatusTxn: txnStatusMap[randomTxStatus]}

	ctx.JSON(http.StatusOK, txnStatus)

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
	}

	// if last number >= 5, the card is a customer bank card valid, then
	// response OK with the cardUUID for later communication between
	// the payment gateway and the bank
	if lastNumb >= 5 {
		cardUUID := struct {
			CardUUID string `json:"card_uuid"`
		}{CardUUID: uuid.New().String()}

		ctx.JSON(http.StatusOK, cardUUID)
	} else {
		ctx.JSON(http.StatusUnauthorized, nil)
	}
}
