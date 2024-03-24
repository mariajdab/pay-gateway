package customers

import "github.com/mariajdab/pay-gateway/internal/entity"

type UseCase interface {
	CreateCustomer(c entity.Customer) (string, error)
}
