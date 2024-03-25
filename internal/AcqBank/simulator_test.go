package AcqBank

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckTxnApproved(t *testing.T) {
	r := gin.Default()

	RegisterHTTPEndpoints(r, &Handler{})

	w := httptest.NewRecorder()

	txnInfo := entity.TxnInfo{
		// the cardUUID final char is a letter, so the txn is approved
		CardUUID: "3e2a16c7-002b-4e1a-9f12-75d4bfd9832c",
		Amount:   120.4,
		Currency: "USD",
	}

	body, err := json.Marshal(&txnInfo)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/bank-sim.com/transactions/validate", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	statusTx := &TxnResp{}

	err = json.Unmarshal(w.Body.Bytes(), statusTx)
	assert.Equal(t, err, nil)

	// 2020 status approved
	assert.Equal(t, statusTx.Status, uint16(2020))
}

func TestCheckTxnDeclined(t *testing.T) {
	r := gin.Default()

	RegisterHTTPEndpoints(r, &Handler{})

	w := httptest.NewRecorder()

	txnInfo := entity.TxnInfo{
		// the cardUUID final char is an int and is even, so the txn is declined
		CardUUID: "3e2a16c7-002b-4e1a-9f12-75d4bfd98324",
		Amount:   457.8,
		Currency: "USD",
	}

	body, err := json.Marshal(&txnInfo)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/bank-sim.com/transactions/validate", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	statusTx := &TxnResp{}

	err = json.Unmarshal(w.Body.Bytes(), statusTx)
	assert.Equal(t, err, nil)

	// 9999 status declined
	assert.Equal(t, statusTx.Status, uint16(9999))
}

func TestValidateCardOK(t *testing.T) {
	r := gin.Default()

	RegisterHTTPEndpoints(r, &Handler{})

	card := entity.Card{
		// the final number is >= 5, so the card should be valid
		Number:  "28942958",
		CVV:     "123",
		ExpDate: time.Now().Add(time.Hour * 24 * 10).Format(time.DateOnly),
	}

	ownerCard := entity.Customer{
		FirstName: "luis",
		LastName:  "paez",
		Address:   "calle 1",
	}

	cardBody := entity.CardData{
		CardInfo:  card,
		OwnerInfo: ownerCard,
	}

	body, err := json.Marshal(&cardBody)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/bank-sim.com/cards/validate", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	cardResp := &CardValidResp{}

	err = json.Unmarshal(w.Body.Bytes(), cardResp)
	assert.NoError(t, err)

	// if the card is valid the response should have the cardUUID of len 36
	assert.Equal(t, len(cardResp.CardUUID), 36)

	// valid field should be true
	assert.Equal(t, cardResp.Valid, true)
}

func TestValidateCardInvalid(t *testing.T) {
	r := gin.Default()

	RegisterHTTPEndpoints(r, &Handler{})

	card := entity.Card{
		// the final number is >= 5, so the card should be invalid
		Number:  "28942954",
		CVV:     "123",
		ExpDate: time.Now().Add(time.Hour * 24 * 10).Format(time.DateOnly),
	}

	ownerCard := entity.Customer{
		FirstName: "luis",
		LastName:  "paez",
		Address:   "calle 1",
	}

	cardBody := entity.CardData{
		CardInfo:  card,
		OwnerInfo: ownerCard,
	}

	body, err := json.Marshal(&cardBody)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/bank-sim.com/cards/validate", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	cardResp := &CardValidResp{}

	err = json.Unmarshal(w.Body.Bytes(), cardResp)
	assert.NoError(t, err)

	// if the card is not valid the response should not have the cardUUID
	assert.Equal(t, len(cardResp.CardUUID), 0)

	// valid field should be false
	assert.Equal(t, cardResp.Valid, false)
}
