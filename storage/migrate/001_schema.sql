-- +goose Up

create table users (
    created_at      timestamp       not null default now(),
    updated_at      timestamp       not null default now(),
    id              bigint          not null primary key,
    lang            varchar(2)      not null default 'ru',
    ref             varchar(64)     not null default '',
    balance         bigint          not null default 0,
    last_bonus      timestamp       not null default '2001-09-28 01:00:00',
    perks           varchar(64)[]   not null default array[]::varchar(64)[]
);

create table groups (
    id              bigint          not null primary key,
    state           varchar(32)     not null default 'none',
    message_id      bigint          not null default 0
);

create table bets (
    id              serial          not null primary key,
    user_id         bigint          not null,
    chat_id         bigint          not null default 0,
    amount          bigint          not null default 100,
    sign            char            not null default 'r',
    won             boolean         not null default false,
    done            boolean         not null default false
);

create table contests (
    amount          bigint          not null default 100,
    creator_id      bigint          not null,
    chat_id         bigint          not null default 0,
    id              serial          not null primary key,
    done            boolean         not null default false,
    participants    varchar(256)[]  not null default ARRAY[],
    canceled        boolean         not null default false,
    winner_id       bigint          not null default 0
);

create index idx_bets_user_id ON bets (user_id);