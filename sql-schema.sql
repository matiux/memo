create table events
(
    id int auto_increment primary key,
    uuid        char(36)     not null,
    playhead    int unsigned not null,
    payload     longtext     not null,
    metadata    longtext     not null,
    recorded_on varchar(255)  not null,
    type        varchar(255) not null,
    constraint UNIQ_5387574AD17F50A634B91FA9 unique (uuid, playhead)
);

create table memo(
    id char(36) primary key,
    body longtext     not null,
    created_at varchar(255)  not null
);