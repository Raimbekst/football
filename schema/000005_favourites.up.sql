CREATE TABLE IF NOT EXISTS favourites(
    id serial not null  unique ,
    user_id int references users(id) on delete cascade not null,
    building_id int references buildings(id) on delete cascade not null
);