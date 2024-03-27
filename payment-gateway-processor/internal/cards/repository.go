package cards

type RepositoryCard interface {
	AddCardDetail(cardToken, expDate string, customerID int) error
	CardTokenExists(code string) (bool, error)
}
