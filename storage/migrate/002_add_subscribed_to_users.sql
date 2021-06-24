-- +goose Up

alter table users add column subscribed bool not null default false;
alter table users drop column lang;
