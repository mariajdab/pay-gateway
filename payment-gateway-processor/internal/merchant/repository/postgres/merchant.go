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

func NewMerchantRepo(conn *pgxpool.Pool) merchant.RepositoryMerchant {
	return &merchantRepo{conn: conn}
}

func (m *merchantRepo) AddMerchant(merchant entity.Merchant) error {
	ctx := context.Background()

	_, err := m.conn.Exec(ctx, `
		INSERT INTO merchants (name, code, account)
		VALUES ($1, $2, $3)
		`,
		merchant.Name,
		merchant.Code,
		merchant.Account)

	return err
}

func (m *merchantRepo) MerchantAccountByCode(code string) (string, error) {
	var account string

	ctx := context.Background()

	if err := m.conn.QueryRow(ctx, `
		SELECT account FROM merchants WHERE code = $1
		`, code).Scan(&account); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return account, nil
}
