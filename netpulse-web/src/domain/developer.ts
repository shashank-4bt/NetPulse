export type ObservationView = {
  value: number | null;
  measured: boolean;
  sampleCount: number;
  summary: string;
};

export type PercentileView = {
  p50: number | null;
  p95: number | null;
  p99: number | null;
  sampleCount: number;
  summary: string;
};

export type RegionalView = {
  region: string;
  sampleCount: number;
  status: string;
  summary: string;
};

export type DeveloperIncidentView = {
  id: string;
  monitorId: string;
  title: string;
  status: string;
  startedAt: string;
  resolvedAt: string | null;
  sampleCount: number;
  summary: string;
};

export type DeveloperDashboardView = {
  availability: ObservationView;
  latency: PercentileView;
  incidents: DeveloperIncidentView[];
  regionalPerformance: RegionalView[];
  summary: string;
};

export type SLAView = {
  availability: ObservationView;
  downtime: ObservationView;
  latency: PercentileView;
  incidents: DeveloperIncidentView[];
  regionalPerformance: RegionalView[];
  summary: string;
};

export type MonitorView = {
  id: string;
  name: string;
  target: string;
  type: string;
  regions: string[];
  frequencySeconds: number;
  timeoutSeconds: number;
  status: string;
  summary: string;
  checkCount: number;
  createdAt: string;
  updatedAt: string;
};

export type MonitorCheckView = {
  id: string;
  monitorId: string;
  region: string;
  ok: boolean;
  latencyMs: number | null;
  at: string;
  summary: string;
};

export type APIKeyView = {
  id: string;
  name: string;
  prefix: string;
  last4: string;
  scopes: string[];
  rateLimitPerMin: number;
  revoked: boolean;
  createdAt: string;
  lastUsedAt: string | null;
};

export type WebhookView = {
  id: string;
  url: string;
  events: string[];
  secretHint: string;
  createdAt: string;
  disabled: boolean;
};

export type DeliveryView = {
  id: string;
  webhookId: string;
  event: string;
  eventId: string;
  timestamp: string;
  idempotencyKey: string;
  signature: string;
  attempt: number;
  status: string;
  summary: string;
};

export type AlertRuleView = {
  id: string;
  kind: string;
  monitorId: string | null;
  threshold: number;
  enabled: boolean;
  deliveredCount: number;
  summary: string;
};

export type UsageView = {
  requests: number;
  measurements: number;
  monitors: number;
  regions: number;
  webhooks: number;
  summary: string;
};

export const DEVELOPER_NAV = [
  { href: "/developers/dashboard", label: "Dashboard" },
  { href: "/developers/monitors", label: "Monitors" },
  { href: "/developers/incidents", label: "Incidents" },
  { href: "/developers/api", label: "API keys" },
  { href: "/developers/webhooks", label: "Webhooks" },
  { href: "/developers/usage", label: "Usage" },
  { href: "/developers/sla", label: "SLA" },
] as const;

export const WEBHOOK_EVENTS = [
  "incident.created",
  "incident.updated",
  "incident.resolved",
  "monitor.down",
  "monitor.recovered",
  "threshold.exceeded",
  "diagnosis.completed",
] as const;

export const emptyObservation = (topic: string): ObservationView => ({
  value: null,
  measured: false,
  sampleCount: 0,
  summary: `Observed sample count: 0. ${topic} is not estimated for a population.`,
});

export const emptyPercentiles = (): PercentileView => ({
  p50: null,
  p95: null,
  p99: null,
  sampleCount: 0,
  summary: "Observed sample count: 0. Percentiles are not estimated from an empty series.",
});

export const emptyDeveloperDashboard = (): DeveloperDashboardView => ({
  availability: emptyObservation("Availability"),
  latency: emptyPercentiles(),
  incidents: [],
  regionalPerformance: [],
  summary: "No monitor checks are stored. Availability and latency stay unmeasured.",
});

export const emptySLA = (): SLAView => ({
  availability: emptyObservation("Availability"),
  downtime: emptyObservation("Downtime"),
  latency: emptyPercentiles(),
  incidents: [],
  regionalPerformance: [],
  summary: "No SLA window can be computed until monitor checks are stored.",
});
