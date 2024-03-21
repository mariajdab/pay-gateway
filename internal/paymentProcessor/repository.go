package paymentProcessor

import "github.com/mariajdab/pay-gateway/internal/entity"

type Repository interface {
	CreateCustomerProcessor(customer entity.Customer)
	GetProcessorIDByCustomerEmail(email string) (id string)
}
