BEGIN;
create table if not exists groups
(

    id   integer generated always as identity
        primary key,
    year integer not null,
    spec text    not null,
    num  integer not null
);

create table if not exists chats
(
    id         bigint generated always as identity
        primary key,
    chat_id    bigint                                 not null
        unique,
    group_id   integer
        references groups,
    created_at timestamp with time zone default now() not null,
    deleted_at timestamp with time zone
);

create table if not exists crons
(
    id      bigint generated always as identity
        primary key,
    chat_id bigint       not null
        references chats,
    title   varchar(200) not null,
    send_at time         not null,
    flags   smallint     not null

);
COMMIT;
