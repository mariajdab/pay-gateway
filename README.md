# Pay-gateway

### Overview

This project attempts to create an online payment platform to help e-commerce businesses. To show a 'real' process, we will have two servers opened: one for bank simulation, related to bank responses, for example, to create a new transaction for card/account or validate card information; the other server is related to the payment platform. It should be able to process payments, retrieve past payments linked to a txnUUID, and make refunds. 

### Desing

- DB Desing

![tables](https://github.com/mariajdab/pay-gateway/blob/main/database-diagrams.svg)


### API Endpoints Payment Platform 

###### Process a Payment: 
##### [POST] http://localhost:8080/processor-pay/payments

Example of a valid request body: 
```
{
  "billing_amount": 150,
  "currency": "USD",
  "card_info": {
    "number": "377673221487787",
    "cvv": "110",
    "exp_date": "2025-03-26"
  },
  "crated_at": "2024-03-26T21:50:01.000Z",
  "merchant_code": "1234#",
  "customer_data": {
    "first_name": "lula",
    "last_name": "Rodriguez",
    "email": "ma@ula.ve",
    "address": "las amaricas",
    "country": "MEX"
  }
}
```

### Note: 

The data should be "almost valid" 

  1. The `card number` format is validated (you could use a card number generator online), `cvv` have length equal to 3, `expiration card` should be a date 

  2. The customer or card owner information is also validated: `email` should have a valid format, `first name` and `last name` with a max of 10 and 12, `address` max leght of 18 and `country` shoud have alpha3Code format, for example: MEX, COL, BRA, etc

  3. At the moment the payment platform only accept USD as `currency` and with a minimum amount of `0.01`

Response example: 

```
{
  "status_payment": "Success",
  "txn_uuid": "6fc49459-81a3-4d00-a163-1356171cf10e"
}
```


###### Retrieve a Payment: 
##### [GET] http://localhost:8080/processor-pay/payments/:merchant_code/:txn_uuid

Reponse example: 
```
{
  "billing_amount": 150,
  "status": "Success",
  "currency": "USD",
  "create_at": "2024-03-27T22:15:27.675915Z",
  "customer_data": {
    "first_name": "lula",
    "last_name": "Rodriguez",
    "email": "ma@ula.ve",
    "address": "las amaricas",
    "country": "MEX"
  }
}
```


###### Refund a Payment: 
##### [POST] http://localhost:8080/processor-pay/payments/refund

Request body example
```
{
"amount": 50,
"txn_uuid": "6fc49459-81a3-4d00-a163-1356171cf10e",
"merchant_code": "1234#"

}
```

Response example: 
```
{
  "status_refund": "Rejected",
  "reason": "the txn is already refunded"
}
```

<br>

#### Bank Simulator:

It handles the simulation response from a bank and not the real logic. For that it have a server active with two endpoints to be called for 
the payment gateway. One endpoint related to validate a card information and the other to try to create a transaction. 
The simulation code tried to be "deterministic", for example to create a transaction the bank code only check the amount request and if it is 
an even number the txn is declined by the "bank", similar case happen for the refund response. In the other hand the validate card
handler use another "simulation logic", and it depends on the `gatwy_pamt_uuid` that is created before send the validation card request to the bank,
so the card will be invalid in the case the uuid ends with a number and valid when is a letter. More detail could be found the file `simulator_test.go` 

## Starting 🚀

git clone https://github.com/mariajdab/pay-gateway.git

### Open a new terminal window 
Copy and run the command ```docker-compose up --build```

#### What happened!? 🚀
The merchants table needs some entries so the `init.sql` will run to add 2 merchants to the db in order to be able to use the api endpoints 
```
INSERT INTO merchants (name, code, account)
VALUES ('tienda-1', '1234#', 'sjlgjljsg934t93tial');
INSERT INTO merchants (name, code, account)
VALUES ('levis', '33342#', '242598fjslflj9320xd');
````
The server for the payment gateway will be enabled and also the server that represents the bank simulation



