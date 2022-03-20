CREATE TABLE IF NOT EXISTS buildings(
    id serial not null unique ,
    building_name varchar(255) not null,
    address varchar(255) not null,
    instagram varchar(255) not null default '',
    manager_id int references users(id) on delete cascade not null,
    description text not null default '',
    work_time int check ( buildings.work_time >= 1 and 2 <= buildings.work_time),
    start_time time not null default now(),
    end_time time not null default now()
);
