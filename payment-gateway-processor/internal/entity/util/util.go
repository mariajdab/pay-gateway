package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/mariajdab/pay-gateway/internal/entity"
	"log"
)

// Sample public key provided by the bank
var rsaBytes = `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEApW1W9dnfdFF7FHrq6HPveR/9T+nM70yO7QOGytR0j/chMBJcJBjG
hJOuKPFbkVyS+BE/4M8CojLgvz4ex82Re0sFa5TqnoWvuP5P4vktR6M5W53sTW3y
gUnfF/oHcEmARQ1xKZdgVnlIfrdbpjecPyLi1Ng4HmhEfCFUOW64koxpb4XeH5O5
q+vc/731ExVOYBU8Sl6kPdjpJuVjS3DHKAVgfVEhscXd3JDjDuMDT3w1IYNb5c2s
wHE55q4Jnc1cr42jdynnkXzmuOGo2C6yD95kbBDLp7wSiBxaMA8gbRkzWJ99T+6l
KsKG2zfndMF3jZW1v1wWiEbYRN07qbN0NQIDAQAB
-----END RSA PUBLIC KEY-----
`

func EncryptCardData(card entity.CardData) ([]byte, error) {
	block, _ := pem.Decode([]byte(rsaBytes))

	if block == nil {
		log.Println("Failed to decode PEM block containing public key")
		return nil, errors.New("error during decode PEM block")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		log.Println("Failed parsing public key")
		return nil, errors.New("error during ParsePKIXPublicKey")
	}

	cardDataByte := []byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s",
		card.CardInfo.ExpDate,
		card.CardInfo.Number,
		card.OwnerInfo.Country,
		card.CardInfo.CVV,
		card.OwnerInfo.Country,
		card.OwnerInfo.FirstName,
		card.OwnerInfo.LastName,
	))

	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, cardDataByte)
}
