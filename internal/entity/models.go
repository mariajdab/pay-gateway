package entity

import "time"

type Customer struct {
	FirstName string
	LastName  string
	Email     string
	Address   string
	Country   string
}

type Card struct {
	Number  string // this will be saved as a token
	CVV     uint16
	ExpDate time.Time
}

type Merchant struct {
	Name string
	Code string
}

type PaymentRequest struct {
	BillingAmount float32
	Currency      string
	Country       string
	CardInfo      Card
	CratedAt      time.Time
	CustomerData  Customer
}

type Transaction struct {
	TxnUUID       string
	BillingAmount float32
	CreatedAt     time.Time
	Status        string
	Currency      string
	CardTk        string
	MerchantCode  string
}

type PaymentInfo struct {
	BillingAmount float32
	Status        string
	Currency      string
	CreateAt      time.Time
	CustomerData  Customer
}
