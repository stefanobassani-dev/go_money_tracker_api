drop table if exists tink_credentials;
drop table if exists users;
drop table if exists tink_webhooks;
drop index if exists idx_users_tink_id;
drop index if exists idx_creds_user_id;
-- 1. Estensione per gli UUID (se non attiva)
CREATE
EXTENSION IF NOT EXISTS "pgcrypto";

-- 2. Tabella Utenti Principale
CREATE TABLE IF NOT EXISTS users
(
    user_id
    UUID
    PRIMARY
    KEY
    DEFAULT
    gen_random_uuid
(
),
    email TEXT UNIQUE NOT NULL,

    -- Identificativo univoco dell'utente su Tink
    tink_user_id TEXT UNIQUE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP
                         WITH TIME ZONE DEFAULT NOW()
    );

-- 3. Tabella Credenziali (Banche collegate)
-- Un utente può avere più banche (es. UniCredit e Intesa)
CREATE TABLE IF NOT EXISTS tink_credentials
(
    id
    UUID
    PRIMARY
    KEY
    DEFAULT
    gen_random_uuid
(
),
    user_id UUID NOT NULL REFERENCES users
(
    user_id
) ON DELETE CASCADE,

    -- ID della specifica connessione bancaria fornito da Tink
    credentials_id TEXT UNIQUE NOT NULL,

    -- Provider (es: "it-unicredit-ob")
    provider_name TEXT,

    -- Stato (es: "UPDATED", "AUTHENTICATING")
    status TEXT,

    created_at TIMESTAMP
  WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP
  WITH TIME ZONE DEFAULT NOW()
    );

CREATE TABLE tink_webhooks
(
    webhook_id                  UUID PRIMARY KEY,                     -- Usiamo l'ID che ci restituisce Tink
    description         TEXT,
    url                 TEXT        NOT NULL,
    secret              TEXT        NOT NULL,                 -- Fondamentale per validare la firma X-Tink-Signature
    disabled            BOOLEAN     DEFAULT false,
    enabled_events      JSONB       NOT NULL,                 -- Salviamo l'array ["refresh:finished", ...]
    created_at          TIMESTAMPTZ NOT NULL,
    updated_at          TIMESTAMPTZ NOT NULL,
    internal_created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP -- Timestamp del nostro DB
);

-- Indici per performance
CREATE INDEX IF NOT EXISTS idx_users_tink_id ON users (tink_user_id);
CREATE INDEX IF NOT EXISTS idx_creds_user_id ON tink_credentials (user_id);