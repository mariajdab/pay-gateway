package customers

import "github.com/mariajdab/pay-gateway/internal/entity"

type UseCase interface {
	CreatePaymentRequest(txn entity.Transaction)
}
