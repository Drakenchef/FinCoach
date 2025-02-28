CREATE TABLE users
(
    id serial not null unique,
    login varchar(255) not null unique,
    password varchar(255) not null
);

CREATE TABLE planets
(
    id serial not null unique,
    name varchar(255) not null,
    description TEXT,
    radius float not null,
    distance float not null,
    gravity float not null,
    image varchar(255),
    is_delete BOOLEAN
);

CREATE TABLE flight_requests (
    id serial not null unique,
    date_start DATE,
    date_end DATE,
    status VARCHAR(255),
    AMS VARCHAR(255),
    user_id INT REFERENCES users(id)
);

CREATE TABLE planets_requests(
    id serial not null unique,
    fr_id INT REFERENCES flight_requests(id),
    planet_id INT REFERENCES planets(id)
)