package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/mariajdab/pay-gateway/internal/cards"
	"github.com/mariajdab/pay-gateway/internal/customers"
	"github.com/mariajdab/pay-gateway/internal/merchant"
	"github.com/mariajdab/pay-gateway/internal/payments"
)

func RegisterHandler(g *gin.Engine,
	paymentUC payments.UseCasePayment,
	cardUC cards.UseCaseCard,
	merchantUC merchant.UseCaseMerchant,
	customerUC customers.UseCaseCustomer) {

	h := NewHandler(paymentUC, cardUC, merchantUC, customerUC)

	proc := g.Group("/processor-pay")
	{
		proc.POST("/payments", h.ProcessPayment)
		proc.GET("/payments/:merchant_code/:payment_uuid", h.RetrievePayment)
		proc.POST("/payments/refund", h.RefundPayment)
	}
}
