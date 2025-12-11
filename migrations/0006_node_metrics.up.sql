CREATE TABLE node_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    
    -- CPU metrics
    cpu_usage_percent DECIMAL(5,2),
    cpu_cores INTEGER,
    
    -- Memory metrics
    memory_total_bytes BIGINT,
    memory_used_bytes BIGINT,
    memory_available_bytes BIGINT,
    memory_usage_percent DECIMAL(5,2),
    
    -- Disk metrics
    disk_total_bytes BIGINT,
    disk_used_bytes BIGINT,
    disk_available_bytes BIGINT,
    disk_usage_percent DECIMAL(5,2),
    
    -- Timestamp
    collected_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Índices para consultas eficientes
CREATE INDEX idx_node_metrics_node_id ON node_metrics(node_id);
CREATE INDEX idx_node_metrics_collected_at ON node_metrics(collected_at DESC);
CREATE INDEX idx_node_metrics_node_time ON node_metrics(node_id, collected_at DESC);

