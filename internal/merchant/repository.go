package merchant

import "github.com/mariajdab/pay-gateway/internal/entity"

type RepositoryMerchant interface {
	AddMerchant(p entity.Merchant) error
	MerchantCodeExists(code string) (bool, error)
}
