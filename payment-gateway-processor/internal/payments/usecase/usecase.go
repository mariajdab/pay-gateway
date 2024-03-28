package usecase

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"time"

	v "github.com/go-ozzo/ozzo-validation"
	"github.com/mariajdab/pay-gateway/internal/cards"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/merchant"
	"github.com/mariajdab/pay-gateway/internal/payments"
)

const (
	currencyUSD = "USD" // currency we accept
)

type paymentUC struct {
	paymentRepo  payments.RepositoryPayment
	cardRepo     cards.RepositoryCard
	merchantRepo merchant.RepositoryMerchant
}

func NewPaymentUC(py payments.RepositoryPayment, card cards.RepositoryCard, merchant merchant.RepositoryMerchant) payments.UseCasePayment {
	return &paymentUC{
		paymentRepo:  py,
		cardRepo:     card,
		merchantRepo: merchant,
	}
}

func (p *paymentUC) ValidatePaymentReq(paymtReq entity.PaymentRequest) (entity.PaymtValidateResp, error) {
	err := v.ValidateStruct(&paymtReq,
		// General payment req validation
		v.Field(&paymtReq.BillingAmount, v.Required, v.Min(0.01)),  // a valid amount must be >= 0.01
		v.Field(&paymtReq.Currency, v.Required, v.In(currencyUSD)), // at the moment we only accept USD as currency
		v.Field(&paymtReq.CratedAt, v.Required),
		v.Field(&paymtReq.MerchantCode, v.Required),
		// Card validation
		v.Field(&paymtReq.CardInfo), // for more detail see the validate method for entity.Card
		// Customer validation
		v.Field(&paymtReq.CustomerData), // for more detail see then validate method for entity.Customer
	)

	if err != nil {
		return entity.PaymtValidateResp{}, err
	}

	merchantAcct, err := p.merchantRepo.MerchantAccountByCode(paymtReq.MerchantCode)
	if err != nil {
		return entity.PaymtValidateResp{}, err
	}

	if merchantAcct == "" {
		log.Println("ValidatePaymentReq: account not found, the merchant is not registered")
		return entity.PaymtValidateResp{
			Status: entity.FailValidationReq,
		}, nil
	}

	// check if the card is expired
	expDateT, _ := time.Parse("2006-01-02", paymtReq.CardInfo.ExpDate)

	if isCardExpired(expDateT) {
		log.Println("ValidatePaymentReq: card is expired")
		return entity.PaymtValidateResp{
			Status: entity.FailValidationReq,
		}, nil
	}

	// get the card "token" from the card number
	cardTk := cardToken(paymtReq.CardInfo.Number)

	// find if exist that card in our system
	existsCard, err := p.cardRepo.CardTokenExists(cardTk)

	if err != nil {
		log.Println("ValidatePaymentReq: error during CardTokenExists")
		return entity.PaymtValidateResp{}, err
	}

	if !existsCard {
		return entity.PaymtValidateResp{
			Status:          entity.PendingBankValidation,
			CardTk:          cardTk,
			MerchantAccount: merchantAcct,
		}, nil
	}

	return entity.PaymtValidateResp{
		Status:          entity.SuccessfulValidation, // we have the card "info" in our system
		CardTk:          cardTk,
		MerchantAccount: merchantAcct,
	}, nil
}

// SavePaymentInfo saves the useful information about the payment
func (p *paymentUC) SavePaymentInfo(txn entity.Transaction) error {
	return p.paymentRepo.AddPaymentTxnHistory(txn)
}

func (p *paymentUC) PaymentDetail(txnUUID, merchantCode string) (entity.PaymentInfo, error) {
	return p.paymentRepo.GetPaymentDetail(txnUUID, merchantCode)
}

func (p *paymentUC) ValidatePaymentRefundReq(txnUUID string, merchantCode string) (entity.RefundResp, error) {
	var refundResp entity.RefundResp

	tnxInfo, status, err := p.paymentRepo.GetPaymentTxnRefundInfo(txnUUID, merchantCode)
	if err != nil {
		return entity.RefundResp{}, err
	}

	switch status {
	case "Refunded":
		refundResp.Valid = false
		refundResp.Reason = entity.Refunded
	case "Success":
		refundResp.Valid = true
		refundResp.TxnInfo = tnxInfo
	case "":
		// txn not found
		refundResp.Valid = false
		refundResp.Reason = entity.NotFound
	}
	return refundResp, nil
}

func (p *paymentUC) UpdatePaymentTxnStatus(txnUUID, status string) error {
	return p.paymentRepo.UpdatePaymentTxnStatus(txnUUID, status)
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
