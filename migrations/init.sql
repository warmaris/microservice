create table exibillia (
    id bigint primary key auto_increment,
    name varchar(255) not null,
    description varchar(255) not null,
    tags json,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp
);

create table looncan (
    id bigint primary key auto_increment,
    name varchar(255) not null,
    value varchar(255) not null,
    parent_id bigint not null,
    parent_type varchar(50) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp
);

create table acaer_versions (
    id int primary key auto_increment,
    version varchar(16) not null
);

insert into acaer_versions (version) values ('v1.0'), ('v1.2'), ('v2.0');

create table acaer (
    id bigint primary key auto_increment,
    name varchar(32) not null,
    version varchar(16) not null references acaer_versions(version)
);

