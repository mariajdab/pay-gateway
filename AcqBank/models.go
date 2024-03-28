package main

type Customer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Country   string `json:"country"`
}

type TxnInfo struct {
	CardInfoEncrypted []byte  `json:"card_info_encrypted"`
	Amount            float32 `json:"amount"`
	Currency          string  `json:"currency"`
	MerchantAccount   string  `json:"merchant_account"`
}
