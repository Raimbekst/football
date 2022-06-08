CREATE TABLE IF NOT EXISTS times(
    id serial not null unique ,
    work_time time,
    building_id int references buildings(id) on delete cascade not null
);