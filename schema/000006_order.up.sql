CREATE TABLE IF NOT EXISTS orders(
    id serial not null unique ,
    pitch_id int references pitches(id) not null,
    user_id int references users(id) not null,
    order_date timestamp with time zone,
    start_time time not null,
    end_time time not null,
    status int check ( orders.status >= 1 and 2 >= orders.status )
);