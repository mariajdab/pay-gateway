package cards

type Repository interface {
	AddCardDetail(cardToken, cardBankUUID, expDate string, customerID uint32) error
	CardInfoExists(cardTk string) (bool, error)
	GetCardBankUUID(cardUUIDtk string) (string, error)
}
