CREATE TABLE IF NOT EXISTS pitches(
    id serial not null unique ,
    building_id int references buildings(id) on delete cascade not null,
    price int not null ,
    pitch_type int check (pitches.pitch_type >= 1  and 3 <= pitches.pitch_type),
    pitch_extra int check (pitches.pitch_extra >= 1 and 2 <= pitches.pitch_extra),
    pitch_image text not null
);