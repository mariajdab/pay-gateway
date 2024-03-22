package entity

import "time"

type Customer struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type Card struct {
	Number  string    `json:"card_number"`
	CVV     uint16    `json:"cvv"`
	ExpDate time.Time `json:"exp_date"`
	Address string    `json:"address"`
}

type Merchant struct {
	MerchantName string `json:"merchant_name"`
	MerchantCode string `json:"merchant_code"`
}

type PaymentRequest struct {
	BillingAmount float32
	Currency      string
	Country       string
	CardInfo      Card
	CratedAt      time.Time
	MerchantCode  string
	CustomerData  Customer
}

type Transaction struct {
	TxnUUID      string
	UpdatedAt    time.Time
	MerchantCode string
	Status       string
	PaymentReq   PaymentRequest
}

type Account struct {
	ID         int
	CustomerID int
	Currency   string
	Balance    float32
	Status     string
	CardUUID   string
}

type PurchasePayment struct {
	BillingAmount float32
	Currency      string
	Country       string
	Time          time.Time
	CustomerData  Customer
}
