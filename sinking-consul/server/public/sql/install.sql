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
    "group"     varchar(50)  not null,
    name        varchar(100) not null,
    type        varchar(50)  not null,
    hash        varchar(50),
    content     TEXT,
    create_time TEXT,
    update_time TEXT,
    constraint pk_group_name
        primary key ("group", name)
);

create index idx_configs_group
    on cloud_configs ("group");

create index idx_configs_name
    on cloud_configs (name);

create index idx_configs_type
    on cloud_configs (type);

create index idx_configs_hash
    on cloud_configs (hash);

create table cloud_nodes
(
    "group"       varchar(50)       not null,
    name          varchar(50)       not null,
    address       varchar(200)      not null,
    online_status integer default 0 not null,
    status        integer default 0 not null,
    last_heart    integer default 0 not null,
    create_time   TEXT,
    update_time   TEXT,
    constraint pk_group_name_address
        primary key ("group", name, address)
);

create index idx_nodes_group_name
    on cloud_nodes ("group", name);

create index idx_nodes_online_status
    on cloud_nodes (online_status);

create index idx_nodes_status
    on cloud_nodes (status);

create index idx_nodes_last_heart
    on cloud_nodes (last_heart);

