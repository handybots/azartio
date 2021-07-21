create table ideas (
    id          serial              primary key,
    title       varchar(128)        not null,
    description text                not null,
    emoji       varchar(8)          not null,
    used        bool                not null default false,
    deleted     bool                not null default false
);

create table votes (
    created_at  timestamptz         not null default now(),
    updated_at  timestamptz         not null default now(),
    id          serial              primary key,
    message_id  varchar(64)         not null default '',
    done        bool                not null default false,
    days_left   int                 not null,
    ideas       int[]               not null
);

create table voters (
    created_at  timestamptz         not null default now(),
    user_id     bigint              not null,
    vote_id     int                 not null,
    idea_id     int                 not null,

    primary key (user_id, vote_id)
);