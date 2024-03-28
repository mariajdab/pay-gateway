package entity

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// payment request statuses
	PendingBankValidation = "request-bank-to-confirm-card"
	SuccessfulValidation  = "pre-authorized-payment-req"
	FailValidationReq     = "pre-authorized-payment-failed"

	// reasons refund req denied
	Refunded       = "the txn is already refunded"
	NotFound       = "the txn refund request is invalid"
	DeclinedRefund = "declined-by-bank"
)

type Customer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Country   string `json:"country"`
}

type Card struct {
	Number  string `json:"number"` // this will be saved as a token
	CVV     string `json:"cvv"`
	ExpDate string `json:"exp_date"`
}

type CardData struct {
	CardInfo  Card     `json:"card_info"`
	OwnerInfo Customer `json:"owner_info"`
}

type Merchant struct {
	Name    string
	Code    string
	Account string
}

type PaymentRequest struct {
	BillingAmount float32   `json:"billing_amount"`
	Currency      string    `json:"currency"`
	CardInfo      Card      `json:"card_info"`
	CratedAt      time.Time `json:"crated_at"`
	MerchantCode  string    `json:"merchant_code"`
	CustomerData  Customer  `json:"customer_data"`
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
	BillingAmount float32   `json:"billing_amount"`
	Status        string    `json:"status"`
	Currency      string    `json:"currency"`
	CreateAt      time.Time `json:"create_at"`
	CustomerData  Customer  `json:"customer_data"`
}

type TxnInfo struct {
	Amount            float32 `json:"amount"`
	Currency          string  `json:"currency"`
	MerchantAccount   string  `json:"merchant_account"`
	CardInfoEncrypted []byte  `json:"card_info_encrypted"` // suppose with encrypted this card info
}

type RefundResp struct {
	TxnInfo TxnInfo
	Valid   bool
	Reason  string
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
	CardTk          string
	MerchantAccount string
	Status          string
}

type TxnResp struct {
	Status int `json:"status_txn"`
}

type CardValidResp struct {
	Valid bool `json:"valid"`
}

type RefundStatus struct {
	Status uint16 `json:"refund_status"`
}
