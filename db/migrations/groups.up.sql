BEGIN;
create table groups
(

    id   smallint generated always as identity
        primary key ,
    year smallint not null,
    spec text     not null,
    num  smallint not null
);

INSERT INTO groups (year, spec, num)
VALUES (19, 'ВЕБ', 1),
       (19, 'ВЕБ', 2),
       (19, 'ИС', 1),
       (19, 'ИС', 2),
       (19, 'КСК', 1),
       (19, 'КСК', 2),
       (19, 'ОИБ', 1),
       (19, 'ОИБ', 2),
       (19, 'П', 1),
       (19, 'П', 2),
       (19, 'П', 3),
       (19, 'ПД', 1),
       (19, 'ПД', 2),
       (19, 'ПД', 3),
       (19, 'ПД', 4),
       (19, 'СА', 1),
       (19, 'СА', 2),
       (20, 'БД', 1),
       (20, 'ВЕБ', 1),
       (20, 'ВЕБ', 2),
       (20, 'ИС', 1),
       (20, 'ОИБ', 1),
       (20, 'ОИБ', 2),
       (20, 'ПО', 1),
       (20, 'ПО', 2),
       (20, 'ПО', 3),
       (20, 'ПО', 4),
       (20, 'УКСК', 1),
       (20, 'УКСК', 2),
       (20, 'Э', 1),
       (20, 'Э', 2),
       (20, 'ЗИО', 1),
       (20, 'ЗИО', 2),
       (20, 'ЗИО', 3),
       (20, 'П', 1),
       (20, 'П', 3),
       (20, 'ПД', 1),
       (20, 'ПД', 2),
       (20, 'ПД', 3),
       (20, 'ПД', 4),
       (20, 'ПД', 5),
       (20, 'ПСА', 1),
       (20, 'ПСА', 2),
       (20, 'ПСА', 3),
       (20, 'СА', 1),
       (21, 'ИС', 1),
       (21, 'ПСА', 1),
       (21, 'ПСА', 2),
       (21, 'ПСА', 3),
       (21, 'ПСА', 4),
       (21, 'ПСА', 5),
       (21, 'ПСА', 6),
       (21, 'ЗИО', 1),
       (21, 'ЗИО', 2),
       (21, 'ЗИО', 3),
       (21, 'ПД', 1),
       (21, 'ПД', 2),
       (21, 'ПД', 3),
       (21, 'СА', 1),
       (21, 'СА', 2),
       (21, 'Э', 1),
       (21, 'Э', 2),
       (21, 'ВЕБ', 1),
       (21, 'ВЕБ', 2),
       (21, 'ВЕБ', 3),
       (21, 'П', 1),
       (21, 'П', 2),
       (21, 'П', 3),
       (21, 'П', 4),
       (21, 'П', 5),
       (21, 'ПО', 1),
       (21, 'ПО', 2),
       (21, 'ПО', 3),
       (21, 'ПО', 4),
       (21, 'БД', 1),
       (21, 'Л', 1),
       (21, 'Л', 2),
       (21, 'ОИБ', 1),
       (21, 'ОИБ', 2),
       (21, 'ОИБ', 3),
       (21, 'УКСК', 1),
       (21, 'УЛ', 1),
       (22, 'ОИБ', 1),
       (22, 'ОИБ', 2),
       (22, 'П', 1),
       (22, 'П', 2),
       (22, 'П', 3),
       (22, 'П', 4),
       (22, 'П', 5),
       (22, 'ПД', 1),
       (22, 'ПД', 2),
       (22, 'ПО', 1),
       (22, 'ПО', 2),
       (22, 'ПО', 3),
       (22, 'ПСА', 1),
       (22, 'ПСА', 2),
       (22, 'ПСА', 3),
       (22, 'ВЕБ', 1),
       (22, 'ВЕБ', 2),
       (22, 'ЗИО', 1),
       (22, 'ЗИО', 2),
       (22, 'ИС', 1),
       (22, 'Э', 1),
       (22, 'Э', 2),
       (22, 'УКСК', 1),
       (22, 'УКСК', 2),
       (22, 'БД', 1),
       (22, 'Л', 1),
       (22, 'СА', 1),
       (22, 'СА', 2);
COMMIT;