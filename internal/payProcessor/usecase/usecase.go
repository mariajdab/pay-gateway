package usecase

import (
	"crypto/sha1"
	"encoding/hex"
	"time"

	v "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/payProcessor"
)

const (
	// currency
	currencyUSD = "USD"

	// payment request statuses
	pendingBankValidation = "request-bank-to-confirm-card"
	sucessfullValidation  = "pre-authorized-payment-req"
	failedValidation      = "payment-req-not-authorized"
	unableValidation      = "unable-to-pre-authorize"
)

type processorUC struct {
	processorRepo payProcessor.Repository
}

func NewProcessorUC(pr payProcessor.Repository) payProcessor.UseCase {
	return &processorUC{processorRepo: pr}
}

func (c *processorUC) ValidatePaymentReq(paymtReq entity.PaymentRequest) (string, error) {

	err := v.ValidateStruct(&paymtReq,
		// Card validation
		v.Field(&paymtReq.CardInfo.Number, v.Required, is.CreditCard), // make a simple card number validation
		v.Field(&paymtReq.CardInfo.CVV, v.Required, v.Min(000), v.Max(999)),
		v.Field(&paymtReq.CardInfo.ExpDate, v.Required),
		// General payment req validation
		v.Field(&paymtReq.BillingAmount, v.Required, v.Min(0.01)),  // a valid amount must be >= 0.01
		v.Field(&paymtReq.Currency, v.Required, v.In(currencyUSD)), // at the moment we only accept USD as currency
		v.Field(&paymtReq.CratedAt, v.Required, v.Date("2006-01-02")),
		v.Field(&paymtReq.Country, v.Required, is.CountryCode3), // countryCode3 eg: VEN, MEX, COL
		v.Field(&paymtReq.MerchantCode, v.Required),
		// Customer validation
		v.Field(&paymtReq.CustomerData.Email, v.Required, is.Email), // make a simple email validation
		v.Field(&paymtReq.CustomerData.FullName, v.Required),
	)

	if err != nil {
		return failedValidation, err
	}

	// check if the card is expired
	expDate := paymtReq.CardInfo.ExpDate
	if isCardExpired(expDate) {
		return failedValidation, err
	}

	// get the card token from the card number
	cardTk := cardToken(paymtReq.CardInfo.Number)

	// check if it is the first payment made with the card in the payment processor
	cardExists, err := c.processorRepo.PaymentCardInfoExists(cardTk)
	if err != nil {
		return unableValidation, err
	}

	if !cardExists {
		return pendingBankValidation, nil
	}

	return sucessfullValidation, nil
}

func (c *processorUC) SaveCustomerCardInfo(card entity.Card) {

}

// isCardExpired checks if the card expiration date is valid
func isCardExpired(expDate time.Time) bool {
	now := time.Now()
	now.Format("2006-01-02")

	maxDateAllowed, _ := time.Parse("2006-01-02", "2050-01-01")

	if expDate.Before(now) || expDate.Equal(now) || expDate.After(maxDateAllowed) {
		return true
	}
	return false
}

// cardToken simulates a "card tokenization system" by the payment processor
func cardToken(cardNumber string) string {
	hash := sha1.New()
	hash.Write([]byte(cardNumber))
	token := hex.EncodeToString(hash.Sum(nil))
	return token
}
