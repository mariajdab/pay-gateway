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
	// the cardUUID final char is a letter, so the txn is approved
	req, _ := http.NewRequest("GET", "/bank-sim.com/cards/3e2a16c7-002b-4e1a-9f12-75d4bfd9832c", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	statusTx := &CheckTxResp{}

	err := json.Unmarshal(w.Body.Bytes(), statusTx)
	assert.Equal(t, err, nil)

	// 2020 status approved
	assert.Equal(t, statusTx.Status, uint16(2020))
}

func TestCheckTxnDeclined(t *testing.T) {
	r := gin.Default()

	RegisterHTTPEndpoints(r, &Handler{})

	w := httptest.NewRecorder()
	// the cardUUID final char is an int and is even, so the txn is declined
	req, _ := http.NewRequest("GET", "/bank-sim.com/cards/3e2a16c7-002b-4e1a-9f12-75d4bfd98322", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	statusTx := &CheckTxResp{}

	err := json.Unmarshal(w.Body.Bytes(), statusTx)
	assert.Equal(t, err, nil)

	// 9999 status declined
	assert.Equal(t, statusTx.Status, uint16(9999))
}

func TestValidateCardOK(t *testing.T) {
	r := gin.Default()

	RegisterHTTPEndpoints(r, &Handler{})

	cardBody := &entity.Card{
		// the final number is >= 5, so the card should be valid
		Number:  "28942958",
		CVV:     123,
		ExpDate: time.Now().Add(time.Hour * 24 * 10),
	}

	body, err := json.Marshal(cardBody)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bank-sim.com/cards/", bytes.NewBuffer(body))
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

func TestValidateCardWrong(t *testing.T) {
	r := gin.Default()

	RegisterHTTPEndpoints(r, &Handler{})

	cardBody := &entity.Card{
		// the final number is >= 5, so the card should be valid
		Number:  "28942954",
		CVV:     123,
		ExpDate: time.Now().Add(time.Hour * 24 * 10),
	}

	body, err := json.Marshal(cardBody)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bank-sim.com/cards/", bytes.NewBuffer(body))
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
