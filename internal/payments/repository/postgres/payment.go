package postgres

import (
	"context"
	"errors"
	"github.com/mariajdab/pay-gateway/internal/payments"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mariajdab/pay-gateway/internal/entity"
)

type paymentRepo struct {
	conn *pgxpool.Pool
}

func NewPaymentRepo(conn *pgxpool.Pool) payments.RepositoryPayment {
	return &paymentRepo{conn: conn}
}

func (p *paymentRepo) AddPaymentTxnHistory(txn entity.Transaction) error {
	ctx := context.Background()

	_, err := p.conn.Exec(ctx, `
		INSERT INTO payment_processor_hist (txn_uuid, amount, currency, card_token, created_at, merchant_code, status_txn)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
			 `,
		txn.TxnUUID, txn.BillingAmount, txn.Currency,
		txn.CardTk, txn.CreatedAt, txn.MerchantCode,
		txn.Status)

	return err
}

func (p *paymentRepo) GetPaymentTxnStatus(txnUUID string) (string, error) {
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

func (p *paymentRepo) UpdatePaymentTxnStatus(txnUUID, status string) error {
	ctx := context.Background()

	_, err := p.conn.Exec(ctx, `
		UPDATE payment_processor_hist SET status_txn = $1, updated_at = now()
		WHERE txn_uuid = $2
			 `, status, txnUUID)

	return err
}

func (p *paymentRepo) GetPaymentDetailByTxUUID(txUUID string) (entity.PaymentInfo, error) {
	ctx := context.Background()
	var payment entity.PaymentInfo

	sqlPaymentDetail := `
		SELECT 
		    pph.amount as billing_amount,
		    pph.currency,
		    pph.created_at,
		    pph.status_tx as status,
		    c.first_name,
		    c.last_name,
		    c.email
		FROM payment_processor_hist pph
		JOIN cards ON pph.card_tk = cards.card_tk
		JOIN customers c ON cards.customer_id = c.id
		WHERE pph.txn_uuid = $1
		LIMIT 1
	`
	if err := p.conn.QueryRow(ctx, sqlPaymentDetail, txUUID).Scan(
		&payment.BillingAmount, &payment.Currency, &payment.CreateAt,
		&payment.Status, &payment.CustomerData.FirstName, &payment.CustomerData.LastName,
		&payment.CustomerData.Email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.PaymentInfo{}, nil
		}
		return entity.PaymentInfo{}, err
	}

	return payment, nil
}
