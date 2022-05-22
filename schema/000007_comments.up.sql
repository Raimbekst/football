CREATE TABLE IF NOT EXISTS comments(
    id serial  not null unique ,
    comment text not null default '',
    user_id int references users(id) on delete cascade not null,
    building_id int references buildings(id) on delete cascade not null
);

CREATE TABLE IF NOT EXISTS grades(
    id serial  not null unique ,
    user_id int references users(id) on delete cascade not null,
    building_id int references buildings(id) on delete cascade not null,
    grade int check ( grades.grade >= 1 and 5 >= grades.grade )
);
