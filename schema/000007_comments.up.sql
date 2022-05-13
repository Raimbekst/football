CREATE TABLE IF NOT EXISTS comments(
    id serial  not null unique ,
    comment text not null default '',
    grade int check ( comments.grade >= 1 and 5 >= comments.grade ),
    user_id int references users(id) on delete cascade not null,
    building_id int references buildings(id) on delete cascade not null
);