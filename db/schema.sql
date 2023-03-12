-- TODO PROPOSAL[CHAT_PK]: delete column ID, use telegram id as PRIMARY KEY

BEGIN;
create table if not exists groups
(

    id   bigint generated always as identity
        primary key,
    year smallint not null,
    spec text     not null,
    num  smallint not null
);

create table if not exists chats
(
    id            bigint generated always as identity
        primary key,
    chat_id       bigint                                 not null
        unique,
    group_id   smallint
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
