CREATE TABLE IF NOT EXISTS services(
   id serial not null unique ,
   service_name varchar(255) not null,
   price int
);

CREATE TABLE IF NOT EXISTS cards(
    id serial not null unique ,
    full_name varchar(255) not null,
    user_id int references users(id) on delete cascade not null,
    cvv int not null,
    full_number varchar(255) not null
);


CREATE TABLE IF NOT EXISTS orders(
    id serial not null unique ,
    first_name varchar(255),
    phone_number varchar(255) not null,
    extra_info  text,
    card_id int references cards(id) not null,
    pitch_id int references pitches(id) not null,
    user_id int references users(id) not null,
    order_date timestamp with time zone,
    end_order_date timestamp with time zone,
    status int check ( orders.status >= 1 and 2 >= orders.status )
);


CREATE TABLE IF NOT EXISTS order_services (
   id serial not null unique ,
   service_id int references services(id),
   order_id int references orders(id) on delete cascade not null
);

CREATE TABLE IF NOT EXISTS order_times (
    id serial not null unique ,
    order_work_time time,
    order_id int references orders(id) on delete cascade not null
);