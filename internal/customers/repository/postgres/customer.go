package postgres

import (
	"context"
	"github.com/mariajdab/pay-gateway/internal/customers"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mariajdab/pay-gateway/internal/entity"
)

type customerRepo struct {
	conn *pgxpool.Pool
}

func NewCustomerRepo(conn *pgxpool.Pool) customers.Repository {
	return &customerRepo{conn: conn}
}

// AddCustomer saves the customer info associated with a card
func (c *customerRepo) AddCustomer(customer entity.Customer) error {
	ctx := context.Background()

	_, err := c.conn.Exec(ctx, `
		INSERT INTO customers (first_name, last_name, email, address, country)
		VALUES ($1, $2, $3, $4)
		`,
		customer.FirstName,
		customer.LastName,
		customer.Email,
		customer.Address)

	return err
}
