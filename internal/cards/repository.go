package cards

type RepositoryCard interface {
	AddCardDetail(cardToken, cardBankUUID, expDate string, customerID int) error
	GetCardBankUUID(cardTk string) (string, error)
}
