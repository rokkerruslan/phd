-- mi: 005.payments.sql

CREATE TABLE payments (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    account_id INTEGER NOT NULL,
    amount     FLOAT NOT NULL,

    FOREIGN KEY (account_id) REFERENCES accounts(id)
)
