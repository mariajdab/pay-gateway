package postgres

import (
	"context"
	"errors"
	"fmt"
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
func (c *cardRepo) AddCardDetail(cardToken, expDate string, customerID int) error {
	ctx := context.Background()

	fmt.Println("TRY SAVIND CARD INFO", cardToken, expDate, customerID)
	_, err := c.conn.Exec(ctx, `
		INSERT INTO cards (card_token, exp_date, customer_id)
		VALUES ($1, $2, $3)
			 `, cardToken, expDate, customerID)

	return err
}

func (c *cardRepo) CardTokenExists(code string) (bool, error) {
	var exists bool

	ctx := context.Background()

	if err := c.conn.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM cards WHERE card_token = $1)
		`, code).Scan(&exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}
