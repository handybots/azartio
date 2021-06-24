-- +goose Up

create table payments (
    created_at      timestamp       not null default now(),
    id              serial          not null primary key,
    user_id         bigint          not null default 0 references users (id),
    target          varchar(64)     not null default '',
    amount          decimal(10, 2)  not null default 0.0,
    profit          decimal(10, 2)  not null default 0.0,
    pay_at          timestamp       null
);