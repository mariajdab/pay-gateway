package customers

import "github.com/mariajdab/pay-gateway/internal/entity"

type RepositoryCustomer interface {
	AddCustomer(customer entity.Customer) error
}
