CREATE TABLE IF NOT EXISTS notifications(
    id serial not null unique ,
    title varchar(255) not null,
    content text not null
);
