package usecase

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/mariajdab/pay-gateway/internal/cards"
	"time"

	v "github.com/go-ozzo/ozzo-validation"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/payments"
)

const (
	// currencies we currently accept
	currencyUSD = "USD"

	// payment request statuses
	pendingBankValidation = "request-bank-to-confirm-card"
	successfulValidation  = "pre-authorized-payment-req"
	failedValidation      = "payment-req-not-authorized"
	unableValidation      = "unable-to-pre-authorize"
)

type paymentUC struct {
	paymentRepo payments.RepositoryPayment
	cardRepo    cards.RepositoryCard
}

func NewpaymentUC(py payments.RepositoryPayment, c cards.RepositoryCard) payments.UseCasePayment {
	return &paymentUC{
		paymentRepo: py,
		cardRepo:    c,
	}
}

func (p *paymentUC) ValidatePaymentReq(paymtReq entity.PaymentRequest) (string, error) {
	if err := v.ValidateStruct(&paymtReq,
		// Card validation
		v.Field(&paymtReq.CardInfo), // for more detail see validate method for entity.Card

		// General payment req validation
		v.Field(&paymtReq.BillingAmount, v.Required, v.Min(0.01)),  // a valid amount must be >= 0.01
		v.Field(&paymtReq.Currency, v.Required, v.In(currencyUSD)), // at the moment we only accept USD as currency
		v.Field(&paymtReq.CratedAt, v.Required),

		// Customer validation
		v.Field(&paymtReq.CustomerData), // for more detail see validate method for entity.Customer
	); err != nil {
		return failedValidation, err
	}

	// check if the card is expired
	expDateT, err := time.Parse("2006-01-02", paymtReq.CardInfo.ExpDate)
	if err != nil {
		return unableValidation, err
	}

	if isCardExpired(expDateT) {
		return failedValidation, nil
	}

	// get the card token from the card number
	cardTk := cardToken(paymtReq.CardInfo.Number)

	// check if it is the first payment made with the card in the payment processor
	cardExists, err := p.cardRepo.CardInfoExists(cardTk)
	if err != nil {
		return unableValidation, err
	}

	if !cardExists {
		return pendingBankValidation, nil
	}

	return successfulValidation, nil
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
