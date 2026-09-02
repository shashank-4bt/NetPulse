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
