package cards

type RepositoryCard interface {
	AddCardDetail(cardToken, cardBankUUID, expDate string, customerID uint32) error
	CardInfoExists(cardTk string) (bool, error)
	GetCardBankUUID(cardUUIDtk string) (string, error)
}
