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

func (m *merchantRepo) AddMerchant(p entity.Merchant) error {
	ctx := context.Background()

	_, err := m.conn.Exec(ctx, `
		INSERT INTO merchants (name, code)
		VALUES ($1, $2)
		`,
		p.MerchantName,
		p.MerchantCode)

	return err
}

func (m *merchantRepo) MerchantCodeExists(code string) (bool, error) {
	var exists bool

	ctx := context.Background()

	if err := m.conn.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM merchants WHERE code = $1)
		`, code).Scan(&exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}
