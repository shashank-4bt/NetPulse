-- Operational schema. Domain code does not import this file.
-- Apply when NETPULSE_DATABASE_URL points at PostgreSQL.

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  display_name TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  email_verified BOOLEAN NOT NULL DEFAULT false,
  telemetry_opt_in BOOLEAN NOT NULL DEFAULT false,
  alerts JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL,
  last_seen_at TIMESTAMPTZ NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT false,
  user_agent TEXT,
  ip TEXT,
  label TEXT
);

CREATE TABLE IF NOT EXISTS auth_tokens (
  kind TEXT NOT NULL,
  hash TEXT NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  expires_at TIMESTAMPTZ NOT NULL,
  used BOOLEAN NOT NULL DEFAULT false,
  PRIMARY KEY (kind, hash)
);

CREATE TABLE IF NOT EXISTS security_events (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  kind TEXT NOT NULL,
  at TIMESTAMPTZ NOT NULL,
  summary TEXT NOT NULL,
  ip TEXT
);

CREATE TABLE IF NOT EXISTS organizations (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS services (
  slug TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  category TEXT NOT NULL,
  summary TEXT NOT NULL,
  layers TEXT[] NOT NULL
);

CREATE TABLE IF NOT EXISTS diagnoses (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  target_hostname TEXT NOT NULL,
  target_kind TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  report JSONB
);

CREATE TABLE IF NOT EXISTS reports (
  id UUID PRIMARY KEY,
  diagnosis_id UUID NOT NULL REFERENCES diagnoses(id),
  document JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS report_shares (
  token_hash TEXT PRIMARY KEY,
  diagnosis_id UUID NOT NULL REFERENCES diagnoses(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS saved_services (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  slug TEXT NOT NULL REFERENCES services(slug),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, slug)
);

CREATE TABLE IF NOT EXISTS incidents (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  severity TEXT,
  status TEXT,
  scope TEXT NOT NULL,
  started_at TIMESTAMPTZ NOT NULL,
  last_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  affected_services TEXT[],
  regions TEXT[],
  networks TEXT[],
  evidence JSONB,
  hypotheses JSONB,
  confidence JSONB,
  timeline JSONB,
  sample_count INTEGER NOT NULL DEFAULT 0,
  sample_rate TEXT,
  affected_user_count INTEGER
);

CREATE TABLE IF NOT EXISTS subscriptions (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  service_slug TEXT REFERENCES services(slug),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS configuration (
  key TEXT PRIMARY KEY,
  value JSONB NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS workspaces (
  id UUID PRIMARY KEY,
  owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (owner_id)
);

CREATE TABLE IF NOT EXISTS monitors (
  id UUID PRIMARY KEY,
  workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  target TEXT NOT NULL,
  type TEXT NOT NULL,
  regions TEXT[] NOT NULL DEFAULT '{}',
  frequency_seconds INTEGER NOT NULL,
  timeout_seconds INTEGER NOT NULL,
  thresholds JSONB,
  status TEXT NOT NULL,
  summary TEXT NOT NULL,
  check_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS monitor_checks (
  id UUID PRIMARY KEY,
  monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
  workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
  region TEXT NOT NULL,
  ok BOOLEAN NOT NULL,
  latency_ms INTEGER,
  at TIMESTAMPTZ NOT NULL,
  summary TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS api_keys (
  id UUID PRIMARY KEY,
  workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  prefix TEXT NOT NULL,
  last4 TEXT NOT NULL,
  hash TEXT NOT NULL UNIQUE,
  scopes TEXT[] NOT NULL,
  rate_limit_per_min INTEGER NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL,
  last_used_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS webhooks (
  id UUID PRIMARY KEY,
  workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
  url TEXT NOT NULL,
  events TEXT[] NOT NULL,
  secret TEXT NOT NULL,
  secret_hint TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  disabled BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS webhook_deliveries (
  id UUID PRIMARY KEY,
  webhook_id UUID NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
  workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
  event TEXT NOT NULL,
  event_id TEXT NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL,
  idempotency_key TEXT NOT NULL UNIQUE,
  signature TEXT NOT NULL,
  payload TEXT,
  attempt INTEGER NOT NULL DEFAULT 0,
  status TEXT NOT NULL,
  next_retry_at TIMESTAMPTZ,
  summary TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS alert_rules (
  id UUID PRIMARY KEY,
  workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
  kind TEXT NOT NULL,
  monitor_id UUID REFERENCES monitors(id) ON DELETE CASCADE,
  threshold DOUBLE PRECISION NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT true,
  delivered_count INTEGER NOT NULL DEFAULT 0,
  summary TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS developer_incidents (
  id UUID PRIMARY KEY,
  workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
  monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  status TEXT NOT NULL,
  started_at TIMESTAMPTZ NOT NULL,
  resolved_at TIMESTAMPTZ,
  sample_count INTEGER NOT NULL DEFAULT 0,
  summary TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS usage_counters (
  workspace_id UUID PRIMARY KEY REFERENCES workspaces(id) ON DELETE CASCADE,
  requests INTEGER NOT NULL DEFAULT 0,
  measurements INTEGER NOT NULL DEFAULT 0
);
