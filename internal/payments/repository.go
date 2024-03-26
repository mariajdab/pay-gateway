package payments

import "github.com/mariajdab/pay-gateway/internal/entity"

type RepositoryPayment interface {
	AddPaymentTxnHistory(txn entity.Transaction) error
	GetPaymentTxnStatus(txnUUID string) (string, error)
	UpdatePaymentTxnStatus(txnUUID, status string) error
}
