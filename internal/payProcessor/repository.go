package payProcessor

import "github.com/mariajdab/pay-gateway/internal/entity"

type Repository interface {
	AddPaymentCardDetail(card entity.Card, cardTokenUUID, cardBankUUID string) error
	AddPaymentTxnHistory(txn entity.Transaction) error
	GetPaymentTxnStatus(txnUUID string) (string, error)
	UpdatePaymentTxnStatus(txnUUID, status string) error
	GetPaymentCardBankUUID(cardUUIDtk string) (id string, err error)
}
