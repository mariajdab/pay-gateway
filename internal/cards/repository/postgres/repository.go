package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mariajdab/pay-gateway/internal/cards"
)

type cardRepo struct {
	conn *pgxpool.Pool
}

func NewCardRepo(conn *pgxpool.Pool) cards.RepositoryCard {
	return &cardRepo{conn: conn}
}

// AddCardDetail store only safe information about the payment card. TODO : we can encrypt card expiration date
func (c *cardRepo) AddCardDetail(cardToken, cardBankUUID, expDate string, customerID int) error {
	ctx := context.Background()

	_, err := c.conn.Query(ctx, `
		INSERT INTO cards (card_token, exp_date, card_bank_uuid, customer_id)
		VALUES ($1, $2, $3, $4) returning id
			 `, cardToken, expDate, cardBankUUID, customerID)

	return err
}

func (c *cardRepo) GetCardBankUUID(cardTk string) (string, error) {
	var cardUUIDbank string

	ctx := context.Background()

	if err := c.conn.QueryRow(ctx, `
    	SELECT card_bank_uuid FROM cards WHERE card_token = $1 LIMIT 1
  `, cardTk).Scan(&cardUUIDbank); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return cardUUIDbank, nil
}
