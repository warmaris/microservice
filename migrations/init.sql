create table exibillia (
    id bigint primary key auto_increment,
    name varchar(255) not null,
    description varchar(255) not null,
    tags json,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp
);