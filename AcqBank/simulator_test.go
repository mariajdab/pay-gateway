package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTxnApproved(t *testing.T) {
	r := gin.Default()

	RegisterHTTPEndpoints(r)

	w := httptest.NewRecorder()

	txnInfo := TxnInfo{
		MerchantAccount:   "dc2241q",
		CardInfoEncrypted: []byte{'b', '&'},
		// the amount is even so the txn is approved
		Amount:   120.4,
		Currency: "USD",
	}

	body, err := json.Marshal(&txnInfo)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/bank-sim/transactions", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	statusTx := &TxnResp{}

	err = json.Unmarshal(w.Body.Bytes(), statusTx)
	assert.Equal(t, err, nil)

	// 2020 status approved
	assert.Equal(t, statusTx.Status, 2020)
}

func TestCreateTxnDeclined(t *testing.T) {
	r := gin.Default()

	RegisterHTTPEndpoints(r)

	w := httptest.NewRecorder()

	txnInfo := TxnInfo{
		MerchantAccount:   "1480b42",
		CardInfoEncrypted: []byte{'b', '&'},
		// the amount is odd so the txn is declined
		Amount:   457.8,
		Currency: "USD",
	}

	body, err := json.Marshal(&txnInfo)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/bank-sim/transactions", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	statusTx := &TxnResp{}

	err = json.Unmarshal(w.Body.Bytes(), statusTx)
	assert.Equal(t, err, nil)

	// 9999 status declined
	assert.Equal(t, statusTx.Status, 9999)
}

func TestValidateCardOK(t *testing.T) {
	r := gin.Default()

	RegisterHTTPEndpoints(r)

	cq := CardValidateReq{
		CardInfo:      []byte{'b', '&'}, // encrypted card info
		GatwyPamtUUID: "ksjg9258035083amja1",
	}

	body, err := json.Marshal(&cq)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/bank-sim/cards/validate", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	cardResp := &CardValidResp{}

	err = json.Unmarshal(w.Body.Bytes(), cardResp)
	assert.NoError(t, err)

	// the last char of the GatwyPamtUUID is a number so the card should be true
	assert.Equal(t, cardResp.Valid, false)
}

func TestValidateCardInvalid(t *testing.T) {
	r := gin.Default()

	RegisterHTTPEndpoints(r)

	cq := CardValidateReq{
		CardInfo:      []byte{'b', '&'},      // encrypted card info
		GatwyPamtUUID: "ksjg9258035083amjdd", // the last char is a letter, so this should be a txn valid
	}

	body, err := json.Marshal(&cq)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/bank-sim/cards/validate", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	cardResp := &CardValidResp{}

	err = json.Unmarshal(w.Body.Bytes(), cardResp)
	assert.NoError(t, err)

	// valid field should be false
	assert.Equal(t, cardResp.Valid, true)
}
