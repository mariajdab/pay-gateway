package usecase

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/mariajdab/pay-gateway/internal/cards"
	"time"

	v "github.com/go-ozzo/ozzo-validation"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/payments"
)

const currencyUSD = "USD" // currency we accept

type paymentUC struct {
	paymentRepo payments.RepositoryPayment
	cardRepo    cards.RepositoryCard
}

func NewPaymentUC(py payments.RepositoryPayment, card cards.RepositoryCard) payments.UseCasePayment {
	return &paymentUC{
		paymentRepo: py,
		cardRepo:    card,
	}
}

func (p *paymentUC) ValidatePaymentReq(paymtReq entity.PaymentRequest) (entity.PaymtValidateResp, error) {
	err := v.ValidateStruct(&paymtReq,
		// General payment req validation
		v.Field(&paymtReq.BillingAmount, v.Required, v.Min(0.01)),  // a valid amount must be >= 0.01
		v.Field(&paymtReq.Currency, v.Required, v.In(currencyUSD)), // at the moment we only accept USD as currency
		v.Field(&paymtReq.CratedAt, v.Required),
		// Card validation
		v.Field(&paymtReq.CardInfo), // for more detail see the validate method for entity.Card
		// Customer validation
		v.Field(&paymtReq.CustomerData), // for more detail see then validate method for entity.Customer
	)

	if err != nil {
		return entity.PaymtValidateResp{}, err
	}

	// check if the card is expired
	expDateT, err := time.Parse("2006-01-02", paymtReq.CardInfo.ExpDate)
	if err != nil {
		return entity.PaymtValidateResp{}, err
	}

	if isCardExpired(expDateT) {
		return entity.PaymtValidateResp{}, nil
	}

	// get the card token from the card number
	cardTk := cardToken(paymtReq.CardInfo.Number)

	// if it's the first payment made with the card in the payment processor
	// then cardUUIDBank is not in the db
	cardUUIDBank, err := p.cardRepo.GetCardBankUUID(cardTk)
	if err != nil {
		return entity.PaymtValidateResp{}, err
	}

	if cardUUIDBank == "" {
		return entity.PaymtValidateResp{
			Status: entity.PendingBankValidation,
			CardTk: cardTk,
		}, nil
	}

	return entity.PaymtValidateResp{
		Status: entity.SuccessfulValidation,
		CardTk: cardTk,
	}, nil
}

// SavePaymentInfo saves the useful information about the payment
func (p *paymentUC) SavePaymentInfo(txn entity.Transaction) (string, error) {
	// create the uuid associated with the payment that could be used by the merchant
	txn.TxnUUID = uuid.New().String()

	if err := p.paymentRepo.AddPaymentTxnHistory(txn); err != nil {
		return "", err
	}

	return txn.TxnUUID, nil
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
