CREATE TABLE agent_heartbeats (
  id           BIGSERIAL PRIMARY KEY,
  agent_id     UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  at           TIMESTAMPTZ NOT NULL DEFAULT now(),
  payload      JSONB NOT NULL                 -- {agent:{version}, net:{...}}
);

CREATE INDEX idx_hb_agent_at ON agent_heartbeats(agent_id, at DESC);

CREATE TABLE node_inventory_snapshots (
  id           BIGSERIAL PRIMARY KEY,
  node_id      UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
  at           TIMESTAMPTZ NOT NULL DEFAULT now(),
  snapshot     JSONB NOT NULL                 -- cpu/mem/disk/docker/...
);

CREATE INDEX idx_inv_node_at ON node_inventory_snapshots(node_id, at DESC);



CREATE TYPE task_status AS ENUM ('queued','ack','running','succeeded','failed','timeout','canceled');

CREATE TABLE tasks (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  node_id       UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
  agent_id      UUID REFERENCES agents(id) ON DELETE SET NULL, -- quien tomó la tarea
  type          TEXT NOT NULL,                                 -- "exec","docker.pull",...
  payload       JSONB NOT NULL,
  status        task_status NOT NULL DEFAULT 'queued',
  retries       INT NOT NULL DEFAULT 0,
  max_retries   INT NOT NULL DEFAULT 3,
  timeout_sec   INT NOT NULL DEFAULT 120 CHECK (timeout_sec BETWEEN 1 AND 86400),
  created_by    UUID REFERENCES users(id) ON DELETE SET NULL,  -- quién la encoló
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  deadline_at   TIMESTAMPTZ                                   -- calculado = created_at + timeout
);

CREATE INDEX idx_tasks_node_status ON tasks(node_id, status);
CREATE INDEX idx_tasks_created_at ON tasks(created_at DESC);

CREATE TABLE task_logs (
  id         BIGSERIAL PRIMARY KEY,
  task_id    UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
  ts         TIMESTAMPTZ NOT NULL DEFAULT now(),
  stream     TEXT NOT NULL CHECK (stream IN ('stdout','stderr','sys')),
  message    TEXT NOT NULL
);

CREATE INDEX idx_task_logs_task_ts ON task_logs(task_id, ts);
