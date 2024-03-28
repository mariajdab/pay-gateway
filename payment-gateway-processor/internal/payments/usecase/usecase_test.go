package usecase

import (
	"github.com/mariajdab/pay-gateway/internal/entity"
	"github.com/mariajdab/pay-gateway/internal/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValidatePaymentReq(t *testing.T) {
	mockCardRepository := new(mocks.RepositoryCard)
	mockPaymentRepository := new(mocks.RepositoryPayment)
	mockMerchantRepository := new(mocks.RepositoryMerchant)

	d := time.Date(2025, 05, 01, 0, 0, 0, 0, time.UTC).Format(time.DateOnly)

	t.Run("payment request success with card already in the system", func(t *testing.T) {

		mockPayReq := entity.PaymentRequest{
			BillingAmount: 100,
			Currency:      "USD",
			CratedAt:      time.Now().UTC(),
			MerchantCode:  "452#FV",
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

		u := NewPaymentUC(mockPaymentRepository, mockCardRepository, mockMerchantRepository)

		mockMerchantRepository.On("MerchantAccountByCode", mockPayReq.MerchantCode).Return("28492u49", nil).Once()

		ctk := cardToken("377673221487787")
		mockCardRepository.On("CardTokenExists", ctk).Return(true, nil).Once()

		validateResp, err := u.ValidatePaymentReq(mockPayReq)
		assert.NoError(t, err)
		assert.Equal(t, validateResp.Status, entity.SuccessfulValidation)
		assert.Equal(t, validateResp.CardTk, ctk)
	})

	t.Run("payment request failed, card number not valid", func(t *testing.T) {

		mockPayReq := entity.PaymentRequest{
			BillingAmount: 50,
			Currency:      "USD",
			CratedAt:      time.Now().UTC(),
			MerchantCode:  "452#FV",
			CardInfo: entity.Card{
				Number:  "3776737", // invalid format card (you can usa a card number generator to try)
				CVV:     "110",
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

		u := NewPaymentUC(mockPaymentRepository, mockCardRepository, mockMerchantRepository)

		_, err := u.ValidatePaymentReq(mockPayReq)
		// wrong credit number format
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a valid credit card number..")
	})

}
