package merchant

import "github.com/mariajdab/pay-gateway/internal/entity"

type Repository interface {
	AddPPurchasePayment(p entity.PurchasePayment, paymtCode string) error
	GetPurchasePayment(paymtCode string) (entity.PurchasePayment, error)
}
