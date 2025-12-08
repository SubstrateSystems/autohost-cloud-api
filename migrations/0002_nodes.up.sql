
CREATE TABLE nodes (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  hostname     TEXT NOT NULL,
  ip_local     INET,
  os           TEXT,
  arch         TEXT,
  version_agent TEXT,
  owner_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  last_seen_at TIMESTAMPTZ DEFAULT now(),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_nodes_hostname ON nodes(hostname);


CREATE UNIQUE INDEX ux_nodes_owner_host
  ON nodes(owner_id, hostname);