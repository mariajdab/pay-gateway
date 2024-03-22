package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/payProcessor"
)

type paymentProcessorRepo struct {
	conn *pgxpool.Pool
}

func NewPaymentProcessor(conn *pgxpool.Pool) payProcessor.Repository {
	return &paymentProcessorRepo{conn: conn}
}

func (p *paymentProcessorRepo) AddPaymentTxnHistory(txn entity.Transaction) error {
	ctx := context.Background()

	_, err := p.conn.Exec(ctx, `
		INSERT INTO payment_processor_hist (txn_uuid, amount, currency, country, card_token, created_at, purc_merchant_code, status_txn)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 `,
		txn.TxnUUID, txn.BillingAmount, txn.Currency,
		txn.Country, txn.CardUUIDtk, txn.CratedAt,
		txn.MerchantCode, txn.Status)

	return err
}

func (p *paymentProcessorRepo) GetPaymentTxnStatus(txnUUID string) (string, error) {
	var status string

	ctx := context.Background()

	if err := p.conn.QueryRow(ctx, `
    SELECT status_txn FROM payment_processor_hist WHERE txn_uuid = $1 LIMIT 1
  `, txnUUID).Scan(&status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return status, nil
}

func (p *paymentProcessorRepo) UpdatePaymentTxnStatus(txnUUID, status string) error {
	ctx := context.Background()

	_, err := p.conn.Exec(ctx, `
		UPDATE payment_processor_hist SET status_txn = $1, updated_at = now()
		WHERE txn_uuid = $2
			 `, status, txnUUID)

	return err
}

func (p *paymentProcessorRepo) GetPaymentCardBankUUID(cardUUIDtk string) (id string, err error) {
	var cardUUIDbank string

	ctx := context.Background()

	if err := p.conn.QueryRow(ctx, `
    SELECT card_bank_uuid FROM card_processor WHERE card_uuid_token = $1 LIMIT 1
  `, cardUUIDtk).Scan(&cardUUIDbank); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return cardUUIDbank, nil
}

func (p *paymentProcessorRepo) AddPaymentCardDetail(card entity.Card, cardTokenUUID, cardBankUUID string) error {
	ctx := context.Background()

	_, err := p.conn.Exec(ctx, `
		INSERT INTO card_processor (card_uuid_token, last4_digits, exp_date, card_bank_uuid, card_address)
		VALUES ($1, $2, $3, $4)
			 `, cardTokenUUID, card.Last4Digits, card.ExpDate, cardBankUUID, card.Address)

	return err
}
