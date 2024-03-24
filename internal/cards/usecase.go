package cards

type UseCase interface {
	SaveCardInfo(cardToken, cardBankUUID, expDate string, customerID uint32) error
}
