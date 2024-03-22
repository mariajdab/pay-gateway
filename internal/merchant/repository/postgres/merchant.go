package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/merchant"
)

type merchantRepo struct {
	conn *pgxpool.Pool
}

func NewMerchantRepo(conn *pgxpool.Pool) merchant.Repository {
	return &merchantRepo{conn: conn}
}

func (m *merchantRepo) AddPPurchasePayment(p entity.PurchasePayment, paymtCode string) error {
	ctx := context.Background()

	// ideally we can hava a customers table with the fullName, email and address info
	_, err := m.conn.Exec(ctx, `
		INSERT INTO merchant_purchase_payments (amount, currency, country, created_at, full_name, email, paymt_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
			 `,
		p.BillingAmount, p.Currency, p.Country,
		p.Time, p.CustomerData.FullName, p.CustomerData.Email,
		paymtCode)

	return err
}

func (m *merchantRepo) GetPurchasePayment(paymtCode string) (entity.PurchasePayment, error) {
	var paymtDetail entity.PurchasePayment

	ctx := context.Background()

	if err := m.conn.QueryRow(ctx, `
    SELECT amount, currency, country, created_at, full_name, email 
    FROM merchant_purchase_payments WHERE paymt_code = $1 LIMIT 1
  `, paymtCode).Scan(
		&paymtDetail.BillingAmount, paymtDetail.Currency,
		paymtDetail.Country, paymtDetail.CustomerData.FullName,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PurchasePayment{}, nil
		}
		return entity.PurchasePayment{}, err
	}

	return paymtDetail, nil
}
