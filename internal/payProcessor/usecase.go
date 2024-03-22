package payProcessor

import "github.com/mariajdab/pay-gateway/internal/entity"

type UseCase interface {
	ValidatePaymentReq(paymtReq entity.PaymentRequest) (string, error)
	SaveCustomerCardInfo(card entity.Card)
	//GetPaymentFinalStatus(txn entity.Transaction)
}
