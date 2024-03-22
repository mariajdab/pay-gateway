package payProcessor

import "github.com/mariajdab/pay-gateway/internal/entity"

type UseCase interface {
	ValidatePaymentReq(txn entity.Transaction)
	SaveCustomerCardInfo(card entity.Card)
	GetPaymentFinalStatus(txn entity.Transaction)
}
