CREATE TABLE agents (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  node_id         UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
  -- identidad y estado
  version         TEXT NOT NULL,                  -- v0.1.3
  status          TEXT NOT NULL DEFAULT 'inactive' CHECK (status IN ('inactive','active','revoked')),
  enrolled_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_seen_at    TIMESTAMPTZ,                    -- del heartbeat
  last_ip_public  INET,
  last_ip_local   INET,
  os              TEXT,
  arch            TEXT,
  labels          JSONB DEFAULT '{}'::jsonb,      -- tags del nodo/agente
  -- seguridad
  auth_mode       TEXT NOT NULL DEFAULT 'jwt' CHECK (auth_mode IN ('jwt','mtls')),
  pubkey_fpr      TEXT,                           -- fingerprint de clave (si usas mTLS)
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_agents_node ON agents(node_id);
CREATE INDEX idx_agents_status ON agents(status);

CREATE UNIQUE INDEX ux_agents_active_per_node
ON agents(node_id)
WHERE status = 'active';


CREATE TABLE agent_tokens (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  agent_id      UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  token_hash    TEXT NOT NULL,          
  expires_at    TIMESTAMPTZ NOT NULL,
  revoked_at    TIMESTAMPTZ,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_agent_tokens_agent ON agent_tokens(agent_id);
CREATE INDEX idx_agent_tokens_valid
  ON agent_tokens(agent_id, expires_at)
  WHERE revoked_at IS NULL;



CREATE TABLE enroll_tokens (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  token_hash    TEXT NOT NULL,                -- one-time, hash
  uses_max      INT  NOT NULL DEFAULT 1 CHECK (uses_max >= 1),
  uses_count    INT  NOT NULL DEFAULT 0 CHECK (uses_count >= 0),
  node_id       UUID REFERENCES nodes(id) ON DELETE CASCADE,  -- opcional si ya existe el nodo
  user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  expires_at    TIMESTAMPTZ NOT NULL,
  consumed_at   TIMESTAMPTZ,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_enroll_user ON enroll_tokens(user_id);
CREATE INDEX idx_enroll_valid ON enroll_tokens(expires_at) WHERE consumed_at IS NULL;
