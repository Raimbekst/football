CREATE type user_type AS ENUM ('user','manager','admin');

CREATE TABLE IF NOT EXISTS users(
    id serial not null unique,
    user_name varchar(255) not null default '',
    phone_number varchar(100) not null default '',
    password varchar(255) default '',
    registered_at TIMESTAMP with time zone default current_timestamp,
    is_activated boolean not null default false,
    user_type user_type
);

CREATE TABLE IF NOT EXISTS sessions(
    id serial not null unique,
    user_id int references users(id) on delete cascade not null,
    refresh_token varchar(500) not null default '',
    expires_at timestamp with time zone default current_timestamp

);