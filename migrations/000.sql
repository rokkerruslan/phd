
CREATE TABLE accounts (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    email TEXT NOT NULL
);

CREATE TABLE events (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    name        TEXT   NOT NULL,
    description TEXT   NOT NULL,
    owner_id    INT    NOT NULL,

    created     TIMESTAMP NOT NULL,
    updated     TIMESTAMP NOT NULL,
    is_deleted  BOOL      NOT NULL DEFAULT FALSE,

    FOREIGN KEY (owner_id) REFERENCES accounts(id)
);

CREATE TABLE timelines (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    event_id INT    NOT NULL,

    start TIMESTAMP NOT NULL,
    "end" TIMESTAMP NOT NULL,

    place TEXT NOT NULL,

    FOREIGN KEY (event_id) REFERENCES events(id)
);

CREATE TABLE offers (
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    event_id   INT NOT NULL,
    account_id INT NOT NULL,

    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,

    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (account_id) REFERENCES accounts(id),

    CONSTRAINT "only_one_offer_to_event_with_one_account" UNIQUE (event_id, account_id)
);
