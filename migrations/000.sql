
create table events (
    id      serial primary key,
    name    text   not null
);

create table timelines (
    id       serial primary key,
    event_id int    not null,

    start timestamp not null,
    "end" timestamp not null,

    foreign key (event_id) references events(id)
);

alter table timelines add column "point" point not null;

alter table events add column created timestamp not null;
alter table events add column updated timestamp not null;

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

alter table events add column owner_id int;

alter table events add foreign key (owner_id) references accounts(id);
