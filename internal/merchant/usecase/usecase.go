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
	return &merchantUC{
		merchantRepo: merchantR,
	}
}

func (muc *merchantUC) CreateMerchant(merc entity.Merchant) error {
	// generate a merchant code
	merchantCode, err := generateRandomStr()
	if err != nil {
		return errors.New("could not generate merchant code")
	}

	merc.MerchantCode = merchantCode

	return muc.merchantRepo.AddMerchant(merc)
}

func (muc *merchantUC) IsValidMerchant(code string) (bool, error) {
	return muc.merchantRepo.MerchantCodeExists(code)
}

func generateRandomStr() (string, error) {
	s := make([]byte, 10)
	_, err := rand.Read(s)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(s), nil
}
