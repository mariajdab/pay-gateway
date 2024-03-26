package usecases

import (
	"github.com/mariajdab/pay-gateway/internal/cards"
)

type cardUC struct {
	cardRepo cards.RepositoryCard
}

func NewCardUC(c cards.RepositoryCard) cards.UseCaseCard {
	return &cardUC{
		cardRepo: c,
	}
}

func (c *cardUC) SaveCardInfo(cardToken, cardBankUUID, expDate string, customerID int) error {
	return c.cardRepo.AddCardDetail(cardToken, cardBankUUID, expDate, customerID)
}
