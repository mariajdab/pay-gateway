package postgres

import (
	"context"
	"github.com/mariajdab/pay-gateway/internal/customers"
	"github.com/mariajdab/pay-gateway/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type customerRepo struct {
	conn *pgxpool.Pool
}

func NewCustomerRepo(conn *pgxpool.Pool) customers.RepositoryCustomer {
	return &customerRepo{conn: conn}
}

// AddCustomer saves the customer info associated with a card
func (c *customerRepo) AddCustomer(customer entity.Customer) (int, error) {
	ctx := context.Background()
	var customerID int

	if err := c.conn.QueryRow(ctx, `
		INSERT INTO customers (first_name, last_name, email, address, country)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
		`,
		customer.FirstName,
		customer.LastName,
		customer.Email,
		customer.Address,
		customer.Country).Scan(&customerID); err != nil {
		return 0, err
	}

	return customerID, nil
}
