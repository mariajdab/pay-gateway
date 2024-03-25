package payments

import "github.com/mariajdab/pay-gateway/internal/entity"

//go:generate mockery --inpackage --testonly --case underscore --name PaymentRepositoryMock
type RepositoryPayment interface {
	AddPaymentTxnHistory(txn entity.Transaction) error
	GetPaymentTxnStatus(txnUUID string) (string, error)
	UpdatePaymentTxnStatus(txnUUID, status string) error
}
