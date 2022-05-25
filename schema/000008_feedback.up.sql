CREATE TABLE IF NOT EXISTS feedbacks(
    id serial not null unique ,
    user_id int references users(id) on delete cascade not null,
    text text not null default ''
);