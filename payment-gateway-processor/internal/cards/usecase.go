package cards

type UseCaseCard interface {
	SaveCardInfo(cardToken, expDate string, customerID int) error
}
