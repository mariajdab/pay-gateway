package customers

import "github.com/mariajdab/pay-gateway/internal/entity"

type Repository interface {
	AddCustomer(customer entity.Customer) error
}
