-- TODO PROPOSAL[GROUP_TYPE]: create table groups (id, year, spec, num)
-- TODO PROPOSAL[GROUP_TYPE]: change column chats.college_group to coll_group_id (FK to groups.id)
-- TODO PROPOSAL[CHAT_PK]: delete column ID, use telegram id as PRIMARY KEY

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
    id            bigint generated always as identity,
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
    chat_id integer      not null
        constraint crons_chats_id_fk
            references chats,
    title varchar(200) not null,
    send_at time     not null,
    flags   smallint not null
);