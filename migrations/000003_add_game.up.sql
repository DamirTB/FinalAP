CREATE TABLE IF NOT EXISTS games (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone not null default NOW(), 
    name text not null,
    price integer not null,
    genres text[] not NULL,
    version integer NOT NULL DEFAULT 1
);