package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/merchant"
)

const (
	// merchant validation"
	allowedMerchant = "merchant-allowed"
	deniedMerchant  = "merchant-not-allowed"
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

// ValidateMerchant validates if the merchant code is authorized to do a payment request
func (muc *merchantUC) ValidateMerchant(merchantCode string) (string, error) {
	merchantExists, err := muc.merchantRepo.MerchantCodeExists(merchantCode)
	if err != nil {
		// internal error of the app
		return deniedMerchant, err
	}

	if !merchantExists {
		return deniedMerchant, nil
	}

	return allowedMerchant, nil

}

func generateRandomStr() (string, error) {
	s := make([]byte, 10)
	_, err := rand.Read(s)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(s), nil
}
