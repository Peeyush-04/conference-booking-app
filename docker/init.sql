create table if not exists bookings (
    id serial primary key,
    first_name text not null,
    last_name text not null,
    email text not null,
    tickets smallint not null
);