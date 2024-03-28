package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/merchant"
)

type merchantUC struct {
	merchantRepo merchant.RepositoryMerchant
}

func NewMerchantUC(merchantR merchant.RepositoryMerchant) merchant.UseCaseMerchant {
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

	merc.Code = merchantCode

	return muc.merchantRepo.AddMerchant(merc)
}

func generateRandomStr() (string, error) {
	s := make([]byte, 10)
	_, err := rand.Read(s)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(s), nil
}
