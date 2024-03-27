package payments

import (
	"github.com/mariajdab/pay-gateway/internal/entity"
)

type RepositoryPayment interface {
	AddPaymentTxnHistory(txn entity.Transaction) error
	GetPaymentTxnRefundInfo(txnUUID, merchantCode string) (entity.TxnInfo, string, error)
	GetPaymentDetail(txUUID, merchantCode string) (entity.PaymentInfo, error)
	UpdatePaymentTxnStatus(txnUUID, status string) error
}
