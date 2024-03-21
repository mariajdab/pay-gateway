package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/paymentProcessor"
)

type paymentProcessorRepo struct {
	Conn *pgxpool.Pool
}

func NewPaymentProcessor(Conn *pgxpool.Pool) paymentProcessor.Repository {
	return &paymentProcessorRepo{}
}

func (p *paymentProcessorRepo) GetProcessorIDByCustomerEmail(email string) (id string) {
	return ""
}

func (p *paymentProcessorRepo) CreateCustomerProcessor(customer entity.Customer) {

}
