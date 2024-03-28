package merchant

import "github.com/mariajdab/pay-gateway/internal/entity"

type RepositoryMerchant interface {
	AddMerchant(p entity.Merchant) error
	MerchantAccountByCode(code string) (string, error)
}
