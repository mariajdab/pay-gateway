package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/merchant"
)

type merchantUC struct {
	merchantRepo merchant.Repository
}

func NewMerchantUC(merchantR merchant.Repository) merchant.UseCase {
	return &merchantUC{merchantRepo: merchantR}
}

func (c *merchantUC) SavePaymentDetail(pyd entity.PurchasePayment) error {
	paymtCode, err := generateRandomStr()

	if err != nil {
		return errors.New("could not generate payment code")
	}

	return c.merchantRepo.AddPPurchasePayment(pyd, paymtCode)
}

func (c *merchantUC) RetrievePaymentDetail(paymtCode string) (entity.PurchasePayment, error) {
	paymentDetail, err := c.RetrievePaymentDetail(paymtCode)

	if err != nil {
		return entity.PurchasePayment{}, err
	}

	return paymentDetail, nil
}

func generateRandomStr() (string, error) {
	s := make([]byte, 10)
	_, err := rand.Read(s)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(s), nil
}
