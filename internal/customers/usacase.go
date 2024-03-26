package customers

import "github.com/mariajdab/pay-gateway/internal/entity"

type UseCaseCustomer interface {
	CreateCustomer(c entity.Customer) (int, error)
}
