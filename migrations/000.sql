
CREATE TABLE events (
    id          SERIAL PRIMARY KEY,
    name        TEXT   NOT NULL,
    description TEXT   NOT NULL,
    owner_id    INT    NOT NULL,

    created     TIMESTAMP NOT NULL,
    updated     TIMESTAMP NOT NULL,

    FOREIGN KEY (owner_id) REFERENCES accounts(id)
);

CREATE TABLE timelines (
    id       serial primary key,
    event_id int    not null,

    start timestamp not null,
    "end" timestamp not null,

    point point not null,

    foreign key (event_id) references events(id)
);

create table accounts (
    id int primary key generated always as identity,

    email text not null
);

create table offers (
    id int primary key generated always as identity,

    event_id   int not null,
    account_id int not null,

    created timestamp not null,
    updated timestamp not null,

    foreign key (event_id) references events(id),
    foreign key (account_id) references accounts(id),

    constraint "only_one_offer_to_event_with_one_account" unique (event_id, account_id)
);
