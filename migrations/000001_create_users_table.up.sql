CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS user_info (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone not null DEFAULT now(),
    fname VARCHAR(255),
    sname VARCHAR(255),
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    user_role VARCHAR(50) DEFAULT 'user',
    balance INT DEFAULT 10000,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);
