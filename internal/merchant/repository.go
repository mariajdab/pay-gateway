package merchant

import "github.com/mariajdab/pay-gateway/internal/entity"

type Repository interface {
	AddMerchant(p entity.Merchant) error
	MerchantCodeExists(code string) (bool, error)
}
