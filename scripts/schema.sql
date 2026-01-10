DROP TABLE IF EXISTS users CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    user_id      UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),

    email        VARCHAR(255) UNIQUE NOT NULL,
    password     TEXT                NOT NULL,

    tink_user_id VARCHAR(255) UNIQUE,

    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_users_email ON users (email);