package merchant

import "github.com/mariajdab/pay-gateway/internal/entity"

type UseCase interface {
	CreateMerchant(merchant entity.Merchant) error
	IsValidMerchant(code string) (bool, error)
}