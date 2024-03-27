package merchant

import (
	"github.com/mariajdab/pay-gateway/internal/entity"
)

type UseCaseMerchant interface {
	CreateMerchant(merchant entity.Merchant) error
}
