package usecase

import (
	"github.com/mariajdab/pay-gateway/internal/customers"
	"github.com/mariajdab/pay-gateway/internal/entity"
)

type customerUC struct {
	customerRepo customers.RepositoryCustomer
}

func NewcustomerUC(cust customers.RepositoryCustomer) customers.UseCaseCustomer {
	return &customerUC{
		customerRepo: cust,
	}
}

func (c *customerUC) CreateCustomer(customer entity.Customer) (int, error) {
	return c.customerRepo.AddCustomer(customer)
}
