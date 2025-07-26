create table cloud_configs
(
    key         varchar(50) not null
        constraint pk_key
            unique,
    value       TEXT,
    update_time TEXT,
    create_time TEXT        not null
);

create table cloud_logs
(
    id          integer           not null
        constraint pk_id
            primary key autoincrement,
    type        integer default 0 not null,
    ip          varchar(50),
    title       varchar(200),
    content     text,
    update_time text,
    create_time text
);

create index cloud_logs_createTime_index
    on cloud_logs (create_time);

create index cloud_logs_ip_index
    on cloud_logs (ip);

create index cloud_logs_type_index
    on cloud_logs (type);

create index cloud_logs_updateTime_index
    on cloud_logs (update_time);