CREATE TABLE IF NOT EXISTS games (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone not null default NOW(), 
    name VARCHAR(255) not null,
    price integer not null,
    CHECK (price >= 0),
    genres text[] not NULL,
    version integer NOT NULL DEFAULT 1
);
