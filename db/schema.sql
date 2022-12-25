create table if not exists chats
(
    id            serial
        primary key,
    chat_id       bigint                                 not null
        unique,
    college_group text,
    created_at    timestamp with time zone default now() not null,
    deleted_at    timestamp with time zone
);

create table if not exists crons
(
    id      bigserial
        primary key,
    chat_id integer      not null
        constraint crons_chats_id_fk
            references chats,
    title varchar(200) not null,
    send_at time     not null,
    flags   smallint not null
);