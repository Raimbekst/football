CREATE TABLE IF NOT EXISTS images(
    id serial not null unique ,
    building_image text not null,
    building_id int references buildings(id) on delete cascade not null
);