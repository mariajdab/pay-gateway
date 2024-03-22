package payProcessor

import "github.com/mariajdab/pay-gateway/internal/entity"

type Repository interface {
	AddPaymentCardDetail(card entity.Card, cardTk, cardBankUUID string) error
	PaymentCardInfoExists(cardTk string) (bool, error)
	AddPaymentTxnHistory(txn entity.Transaction, cardToken string) error
	GetPaymentTxnStatus(txnUUID string) (string, error)
	UpdatePaymentTxnStatus(txnUUID, status string) error
	GetPaymentCardBankUUID(cardTk string) (id string, err error)
}
