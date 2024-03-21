package entity

import "time"

type Customer struct {
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	CardNumber string `json:"card_number"`
}

type Merchant struct {
	MerchantID   int    `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
	MerchantCode string `json:"merchant_code"`
}

type Transaction struct {
	ID                int
	TransactionStatus string
	BillingAmount     float32
	Currency          string
	Country           string
	Time              time.Time
	MerchantData      Merchant
	CustomerData      Customer
}

type Account struct {
	ID              int
	CustomerID      int
	AccountCurrency string
	Balance         float32
	Status          string
}
