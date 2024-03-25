package usecase

import (
	"testing"
	"time"

	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestValidatePaymentReq(t *testing.T) {
	mockCardRepository := new(mocks.RepositoryCard)
	mockPaymentRepository := new(mocks.RepositoryPayment)

	d := time.Date(2025, 05, 01, 0, 0, 0, 0, time.UTC).Format(time.DateOnly)

	t.Run("payment request success: card is already in the system", func(t *testing.T) {

		mockTask := entity.PaymentRequest{
			BillingAmount: 100,
			Currency:      "USD",
			CratedAt:      time.Now().UTC(),
			CardInfo: entity.Card{
				Number:  "377673221487787",
				CVV:     "345",
				ExpDate: d,
			},
			CustomerData: entity.Customer{
				FirstName: "lula",
				LastName:  "Rodriguez",
				Email:     "ma@ula.ve",
				Address:   "las amaricas",
				Country:   "MEX",
			},
		}

		ctk := cardToken("377673221487787")
		mockCardRepository.On("CardInfoExists", ctk).Return(true, nil).Once()

		u := NewpaymentUC(mockPaymentRepository, mockCardRepository)

		validate, err := u.ValidatePaymentReq(mockTask)
		assert.NoError(t, err)
		assert.Equal(t, validate, successfulValidation)

	})

}
