package entity

import "time"

type Customer struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type Card struct {
	Number      string    `json:"card_number"`
	Last4Digits int       `json:"last_4_digits"`
	ExpDate     time.Time `json:"exp_date"`
	Address     string    `json:"address"`
}

type Merchant struct {
	MerchantID   int    `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
	MerchantCode string `json:"merchant_code"`
}

type Transaction struct {
	TxnUUID       string
	BillingAmount float32
	Currency      string
	Country       string
	CardUUIDtk    string
	CratedAt      time.Time
	UpdatedAt     time.Time
	MerchantCode  string
	Status        string
}

type Account struct {
	ID              int
	CustomerID      int
	AccountCurrency string
	Balance         float32
	Status          string
	CardInfo        Card
}

type PurchasePayment struct {
	BillingAmount float32
	Currency      string
	Country       string
	Time          time.Time
	CustomerData  Customer
}
