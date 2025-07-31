create table cloud_clusters
(
    address       varchar(200)      not null
        constraint pk_address
            primary key,
    online_status integer default 0 not null,
    status        integer default 0 not null,
    last_heart    integer default 0,
    create_time   TEXT,
    update_time   TEXT
);

create index idx_cluster_last_heart
    on cloud_clusters (last_heart);

create index idx_cluster_online_status
    on cloud_clusters (online_status);

create index idx_cluster_status
    on cloud_clusters (status);

create table cloud_configs
(
    "group"     varchar(50) not null,
    key         varchar(50) not null,
    value       TEXT,
    create_time TEXT,
    update_time TEXT
);

create index cloud_configs_group_index
    on cloud_configs ("group");

create index cloud_configs_key_index
    on cloud_configs (key);

create table cloud_services
(
    "group"       varchar(50)       not null,
    name          varchar(50)       not null,
    address       varchar(200)      not null,
    online_status integer default 0 not null,
    status        integer default 0 not null,
    last_heart    integer default 0 not null,
    create_time   TEXT,
    update_time   TEXT
);

create index idx_service_group_name
    on cloud_services ("group", name);

create index idx_service_online_status
    on cloud_services (online_status);

create index idx_service_status
    on cloud_services (status);

create index idx_services_last_heart
    on cloud_services (last_heart);

