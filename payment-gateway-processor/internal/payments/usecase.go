package payments

import (
	"github.com/mariajdab/pay-gateway/internal/entity"
)

type UseCasePayment interface {
	ValidatePaymentReq(paymtReq entity.PaymentRequest) (entity.PaymtValidateResp, error)
	SavePaymentInfo(txn entity.Transaction) error
	PaymentDetail(txnUUID, merchantCode string) (entity.PaymentInfo, error)
	ValidatePaymentRefundReq(txnUUID string, merchantCode string) (entity.RefundResp, error)
	UpdatePaymentTxnStatus(txnUUID, status string) error
}
