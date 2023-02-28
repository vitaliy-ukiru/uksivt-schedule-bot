-- TODO PROPOSAL[GROUP_TYPE]: create table groups (id, year, spec, num)
-- TODO PROPOSAL[GROUP_TYPE]: change column chats.college_group to coll_group_id (FK to groups.id)
-- TODO PROPOSAL[CHAT_PK]: delete column ID, use telegram id as PRIMARY KEY
BEGIN;
create table if not exists chats
(
    id            bigint generated always as identity
        primary key,
    chat_id       bigint                                 not null
        unique,
    college_group text,
    created_at    timestamp with time zone default now() not null,
    deleted_at    timestamp with time zone
);

create table if not exists crons
(
    id      bigint generated always as identity
        primary key,
    chat_id integer      not null
        references chats,
    title   varchar(200) not null,
    send_at time         not null,
    flags   smallint     not null
);
COMMIT;