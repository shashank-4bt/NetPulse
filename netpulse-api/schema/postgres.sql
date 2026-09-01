-- Operational schema. Domain code does not import this file.
-- Apply when NETPULSE_DATABASE_URL points at PostgreSQL.

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
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
