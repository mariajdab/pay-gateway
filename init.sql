CREATE TABLE IF NOT EXISTS merchants
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(50)        NOT NULL,
    code    VARCHAR(10) UNIQUE NOT NULL,
    account VARCHAR(24)        NOT NULL
);

CREATE TABLE IF NOT EXISTS customers
(
    id         SERIAL PRIMARY KEY,
    first_name VARCHAR(18)        NOT NULL,
    last_name  VARCHAR(18)        NOT NULL,
    email      VARCHAR(50) UNIQUE NOT NULL,
    country    VARCHAR(5)         NOT NULL,
    address    VARCHAR(18)        NOT NULL
);


CREATE TABLE IF NOT EXISTS cards
(
    id          SERIAL PRIMARY KEY,
    card_token  VARCHAR(50) UNIQUE NOT NULL,
    exp_date    DATE               NOT NULL,
    customer_id integer REFERENCES customers (id)
);


CREATE TABLE payment_processor_hist
(
    id            SERIAL PRIMARY KEY,
    txn_uuid      VARCHAR(36) UNIQUE NOT NULL,
    amount        DOUBLE PRECISION   NOT NULL,
    currency      VARCHAR(8)         NOT NULL,
    card_token    VARCHAR(50) REFERENCES cards (card_token),
    created_at    TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    merchant_code VARCHAR(10) REFERENCES merchants (code),
    status_txn    VARCHAR(10)        NOT NULL,
    updated_at    TIMESTAMP
);


INSERT INTO merchants (name, code, account)
VALUES ('tienda-1', '1234#', 'sjlgjljsg934t93tial');
INSERT INTO merchants (name, code, account)
VALUES ('levis', '33342#', '242598fjslflj9320xd');