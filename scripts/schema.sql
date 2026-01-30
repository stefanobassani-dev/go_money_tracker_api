DROP TABLE IF EXISTS users CASCADE;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS credentials CASCADE;
DROP INDEX IF EXISTS idx_accounts_credential_id;
DROP TABLE IF EXISTS accounts CASCADE;
DROP INDEX IF EXISTS idx_accounts_user_id;

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

CREATE TABLE credentials
(
    credential_id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tink_credential_id  VARCHAR(255) UNIQUE,
    provider_name       VARCHAR(255) NOT NULL,
    type                VARCHAR(50),
    status              VARCHAR(50)  NOT NULL,
    status_payload      TEXT,
    updated             TIMESTAMP WITH TIME ZONE,
    session_expiry_date TIMESTAMP WITH TIME ZONE,
    tink_user_id        VARCHAR(255),

    user_id             UUID         NOT NULL REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE TABLE accounts
(
    id                       VARCHAR(255) PRIMARY KEY,
    user_id                  VARCHAR(255)   NOT NULL,
    credential_id            VARCHAR(255),
    name                     VARCHAR(255)   NOT NULL,
    type                     VARCHAR(50)    NOT NULL,
    balance                  DECIMAL(15, 2) NOT NULL  DEFAULT 0.00,
    currency                 VARCHAR(3)     NOT NULL,
    financial_institution_id VARCHAR(255)   NOT NULL,
    last_refreshed           TIMESTAMP WITH TIME ZONE,
    created_at               TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_accounts_credential_id ON credentials (tink_credential_id);
CREATE INDEX idx_accounts_user_id ON accounts (user_id);