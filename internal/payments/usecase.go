package payments

import "github.com/mariajdab/pay-gateway/internal/entity"

type UseCasePayment interface {
	ValidatePaymentReq(paymtReq entity.PaymentRequest) (string, error)
}
