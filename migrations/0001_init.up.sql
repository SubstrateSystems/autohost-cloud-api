CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto; 

CREATE TABLE users (
  id            BIGSERIAL PRIMARY KEY,
  email         CITEXT UNIQUE NOT NULL,
  name          TEXT,
  password_hash TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- refresh tokens con rotación y revocación
CREATE TABLE refresh_tokens (
  id            BIGSERIAL PRIMARY KEY,
  user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash    TEXT NOT NULL,           -- guarda SOLO el hash del refresh token
  user_agent    TEXT,
  ip            INET,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  revoked_at    TIMESTAMPTZ
);

CREATE INDEX idx_refresh_user ON refresh_tokens(user_id);


-- CREATE TABLE users (
--   id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--   email         CITEXT UNIQUE NOT NULL,
--   name          TEXT,
--   password_hash TEXT NOT NULL,   
--   created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
-- );

-- CREATE TABLE agents (
--   id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--   name         TEXT NOT NULL,
--   labels       JSONB NOT NULL DEFAULT '{}'::jsonb,
--   version      TEXT,
--   os           TEXT,
--   arch         TEXT,
--   last_seen_at TIMESTAMPTZ,
--   created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
-- );
-- CREATE INDEX idx_agents_name ON agents(name);
-- CREATE INDEX idx_agents_labels_gin ON agents USING GIN (labels);

-- CREATE TABLE agent_keys (
--   id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--   agent_id     UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
--   secret_hash  TEXT NOT NULL,
--   created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
--   revoked_at   TIMESTAMPTZ
-- );
-- CREATE INDEX idx_agent_keys_agent ON agent_keys(agent_id);

-- CREATE TABLE api_tokens (
--   id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--   user_id      UUID REFERENCES users(id) ON DELETE SET NULL,
--   name         TEXT NOT NULL,
--   token_hash   TEXT NOT NULL,
--   created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
--   last_used_at TIMESTAMPTZ
-- );
-- CREATE INDEX idx_api_tokens_user ON api_tokens(user_id);
