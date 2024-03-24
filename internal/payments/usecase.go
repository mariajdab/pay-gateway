package payments

import "github.com/mariajdab/pay-gateway/internal/entity"

type UseCase interface {
	ValidatePaymentReq(paymtReq entity.PaymentRequest) (string, error)
	//SavePaymentCardInfo(card entity.Card, cardBankUUID string, customerID uint32) error
	//GetPaymentFinalStatus(txn entity.Transaction)
}
