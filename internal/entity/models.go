package entity

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// merchant validation
	AllowedMerchant = "merchant-allowed"
	DeniedMerchant  = "merchant-not-allowed"

	// payment request statuses
	PendingBankValidation = "request-bank-to-confirm-card"
	SuccessfulValidation  = "pre-authorized-payment-req"
)

type Customer struct {
	FirstName string
	LastName  string
	Email     string
	Address   string
	Country   string
}

type Card struct {
	Number  string // this will be saved as a token
	CVV     string
	ExpDate string
}

type CardData struct {
	CardInfo  Card
	OwnerInfo Customer
}

type Merchant struct {
	Name string
	Code string
}

type PaymentRequest struct {
	BillingAmount float32
	Currency      string
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

type TxnInfo struct {
	CardUUID string
	Amount   float32
	Currency string
}

func (c Card) Validate() error {
	return v.ValidateStruct(&c,
		v.Field(&c.Number, v.Required, is.CreditCard),         // make a simple card number validation
		v.Field(&c.CVV, v.Required, v.Length(3, 3)),           // cvv should be of 3 characters
		v.Field(&c.ExpDate, v.Required, v.Date("2006-01-02")), // should have a format date
	)
}

func (c Customer) Validate() error {
	return v.ValidateStruct(&c,
		v.Field(&c.Email, v.Required, is.Email),            // make a simple email validation
		v.Field(&c.FirstName, v.Required, v.Length(2, 10)), // min and max length limitation
		v.Field(&c.LastName, v.Required, v.Length(2, 12)),  // min and max length limitation
		v.Field(&c.Address, v.Required, v.Length(4, 18)),   // // min and max length limitation
		v.Field(&c.Country, v.Required, is.CountryCode3),   // countryCode3 eg: VEN, MEX, COL
	)
}

type PaymtValidateResp struct {
	CardTk       string
	CardUUIDBank string
	Status       string
}
