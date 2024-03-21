package usecase

import (
	"github.com/mariajdab/pay-gateway/internal/customers"
	"github.com/mariajdab/pay-gateway/internal/customers/repository/postgres"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/paymentProcessor"
)

type CustomerUC struct {
	CustomerRepo  customers.Repository
	ProcessorRepo paymentProcessor.Repository
}

func NewCustomerUC(customerRepo postgres.CustomerRepo) *CustomerUC {
	return &CustomerUC{}
}

func (c *CustomerUC) CreatePaymentRequest(txn entity.Transaction) {
	customerID := c.ProcessorRepo.GetProcessorIDByCustomerEmail(txn.CustomerData.Email)

	if customerID == "" {
		c.ProcessorRepo.CreateCustomerProcessor(txn.CustomerData)
	}
	return
}
