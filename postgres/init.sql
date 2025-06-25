-- User Table
create table if not exists users(
    id serial primary key,
    first_name text not null,
    last_name text not null,
    email text not null unique,
    password_hash text not null,
    role text not null check (role in ('customer', 'organizer')),
    created_at timestamptz not null default now()
);

-- Conference Table
create table if not exists conferences(
    id serial primary key,
    title text not null,
    description text,
    location text not null,
    event_time timestamptz not null,
    total_tickets int not null check(total_tickets > 0),
    available_tickets int not null check(available_tickets >= 0),
    organizer_id int not null references users(id) on delete cascade,
    created_at timestamptz not null default now()
);

-- Booking Table
create table if not exists bookings (
    id serial primary key,
    user_id int not null references users(id) on delete cascade,
    conference_id int not null references conferences(id) on delete cascade,
    tickets_booked int not null check(tickets_booked > 0),
    booked_at timestamptz not null default now()
);

-- Ticket Table
create table if not exists tickets (
    id serial primary key,
    booking_id int not null REFERENCES bookings(id) on delete CASCADE,
    ticket_code text NOT NULL UNIQUE,
    issued_at TIMESTAMPTZ NOT NULL DEFAULT now()
);