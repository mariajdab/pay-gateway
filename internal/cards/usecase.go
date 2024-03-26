package cards

type UseCaseCard interface {
	SaveCardInfo(cardToken, cardBankUUID, expDate string, customerID int) error
}
