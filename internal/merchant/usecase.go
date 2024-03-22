package merchant

import "github.com/mariajdab/pay-gateway/internal/entity"

type UseCase interface {
	SavePaymentDetail(pyd entity.PurchasePayment) error
	RetrievePaymentDetail(paymtCode string) (entity.PurchasePayment, error)
}
