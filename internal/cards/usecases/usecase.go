package usecases

import (
	"github.com/mariajdab/pay-gateway/internal/cards"
)

type cardUC struct {
	cardRepo cards.Repository
}

func NewCardUC(c cards.Repository) cards.UseCase {
	return &cardUC{
		cardRepo: c,
	}
}

func (c *cardUC) SaveCardInfo(cardToken, cardBankUUID, expDate string, customerID uint32) error {
	return c.cardRepo.AddCardDetail(cardToken, cardBankUUID, expDate, customerID)
}
