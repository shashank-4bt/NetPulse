import type {
  APIKeyView,
  AlertRuleView,
  DeliveryView,
  DeveloperDashboardView,
  DeveloperIncidentView,
  MonitorCheckView,
  MonitorView,
  SLAView,
  UsageView,
  WebhookView,
} from "@/domain/developer";
import { emptyDeveloperDashboard, emptyPercentiles, emptySLA } from "@/domain/developer";
import {
  apiRequest,
  asApiFailure,
  type BackendEnvelope,
} from "@/lib/api/backend";
import { type ApiFailure } from "@/lib/api/errors";

function isFail(
  value: { status: number; body: BackendEnvelope } | ApiFailure
): value is ApiFailure {
  return "ok" in value && value.ok === false;
}

async function call(path: string, session: string) {
  return apiRequest(path, { timeoutMs: 10000, session });
}

export async function getDevDashboard(
  session: string
): Promise<{ ok: true; dashboard: DeveloperDashboardView } | ApiFailure> {
  const result = await call("/v1/dev/dashboard", session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.developerDashboard) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, dashboard: mapDashboard(asRecord(result.body.developerDashboard)) };
}

export async function getDevMonitors(
  session: string
): Promise<{ ok: true; monitors: MonitorView[] } | ApiFailure> {
  const result = await call("/v1/dev/monitors", session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, monitors: (result.body.monitors ?? []).map((item) => mapMonitor(asRecord(item))) };
}

export async function getDevMonitor(
  session: string,
  id: string
): Promise<{ ok: true; monitor: MonitorView; checks: MonitorCheckView[] } | ApiFailure> {
  const result = await call(`/v1/dev/monitors/${id}`, session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.monitor) {
    return asApiFailure(result.status, result.body);
  }
  return {
    ok: true,
    monitor: mapMonitor(asRecord(result.body.monitor)),
    checks: (result.body.checks ?? []).map((item) => mapCheck(asRecord(item))),
  };
}

export async function getDevIncidents(
  session: string
): Promise<{ ok: true; incidents: DeveloperIncidentView[] } | ApiFailure> {
  const result = await call("/v1/dev/incidents", session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return {
    ok: true,
    incidents: (result.body.developerIncidents ?? []).map((item) => mapIncident(asRecord(item))),
  };
}

export async function getDevKeys(
  session: string
): Promise<{ ok: true; keys: APIKeyView[] } | ApiFailure> {
  const result = await call("/v1/dev/keys", session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, keys: (result.body.apiKeys ?? []).map((item) => mapKey(asRecord(item))) };
}

export async function getDevWebhooks(
  session: string
): Promise<{ ok: true; webhooks: WebhookView[] } | ApiFailure> {
  const result = await call("/v1/dev/webhooks", session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, webhooks: (result.body.webhooks ?? []).map((item) => mapWebhook(asRecord(item))) };
}

export async function getDevDeliveries(
  session: string,
  id: string
): Promise<{ ok: true; deliveries: DeliveryView[] } | ApiFailure> {
  const result = await call(`/v1/dev/webhooks/${id}/deliveries`, session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return {
    ok: true,
    deliveries: (result.body.deliveries ?? []).map((item) => mapDelivery(asRecord(item))),
  };
}

export async function getDevAlerts(
  session: string
): Promise<{ ok: true; rules: AlertRuleView[] } | ApiFailure> {
  const result = await call("/v1/dev/alerts", session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, rules: (result.body.alertRules ?? []).map((item) => mapAlert(asRecord(item))) };
}

export async function getDevUsage(
  session: string
): Promise<{ ok: true; usage: UsageView } | ApiFailure> {
  const result = await call("/v1/dev/usage", session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.usage) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, usage: mapUsage(asRecord(result.body.usage)) };
}

export async function getDevSLA(
  session: string
): Promise<{ ok: true; sla: SLAView } | ApiFailure> {
  const result = await call("/v1/dev/sla", session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.sla) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, sla: mapSLA(asRecord(result.body.sla)) };
}

export function mapDashboard(value: Record<string, unknown>): DeveloperDashboardView {
  const empty = emptyDeveloperDashboard();
  return {
    availability: mapObservation(asRecord(value.availability), empty.availability),
    latency: mapPercentiles(asRecord(value.latency)),
    incidents: Array.isArray(value.incidents)
      ? value.incidents.map((item) => mapIncident(asRecord(item)))
      : [],
    regionalPerformance: Array.isArray(value.regionalPerformance)
      ? value.regionalPerformance.map((item) => mapRegional(asRecord(item)))
      : [],
    summary: String(value.summary ?? empty.summary),
  };
}

export function mapSLA(value: Record<string, unknown>): SLAView {
  const empty = emptySLA();
  return {
    availability: mapObservation(asRecord(value.availability), empty.availability),
    downtime: mapObservation(asRecord(value.downtime), empty.downtime),
    latency: mapPercentiles(asRecord(value.latency)),
    incidents: Array.isArray(value.incidents)
      ? value.incidents.map((item) => mapIncident(asRecord(item)))
      : [],
    regionalPerformance: Array.isArray(value.regionalPerformance)
      ? value.regionalPerformance.map((item) => mapRegional(asRecord(item)))
      : [],
    summary: String(value.summary ?? empty.summary),
  };
}

function mapObservation(
  value: Record<string, unknown>,
  fallback: { value: number | null; measured: boolean; sampleCount: number; summary: string }
) {
  return {
    value: typeof value.value === "number" ? value.value : null,
    measured: Boolean(value.measured),
    sampleCount: typeof value.sampleCount === "number" ? value.sampleCount : fallback.sampleCount,
    summary: String(value.summary ?? fallback.summary),
  };
}

function mapPercentiles(value: Record<string, unknown>) {
  const empty = emptyPercentiles();
  return {
    p50: typeof value.p50 === "number" ? value.p50 : null,
    p95: typeof value.p95 === "number" ? value.p95 : null,
    p99: typeof value.p99 === "number" ? value.p99 : null,
    sampleCount: typeof value.sampleCount === "number" ? value.sampleCount : 0,
    summary: String(value.summary ?? empty.summary),
  };
}

function mapRegional(value: Record<string, unknown>) {
  return {
    region: String(value.region ?? ""),
    sampleCount: typeof value.sampleCount === "number" ? value.sampleCount : 0,
    status: String(value.status ?? "observed"),
    summary: String(value.summary ?? ""),
  };
}

function mapIncident(value: Record<string, unknown>): DeveloperIncidentView {
  return {
    id: String(value.id ?? ""),
    monitorId: String(value.monitorId ?? ""),
    title: String(value.title ?? ""),
    status: String(value.status ?? ""),
    startedAt: String(value.startedAt ?? ""),
    resolvedAt: typeof value.resolvedAt === "string" ? value.resolvedAt : null,
    sampleCount: typeof value.sampleCount === "number" ? value.sampleCount : 0,
    summary: String(value.summary ?? ""),
  };
}

function mapMonitor(value: Record<string, unknown>): MonitorView {
  return {
    id: String(value.id ?? ""),
    name: String(value.name ?? ""),
    target: String(value.target ?? ""),
    type: String(value.type ?? ""),
    regions: Array.isArray(value.regions) ? value.regions.map((item) => String(item)) : [],
    frequencySeconds: typeof value.frequencySeconds === "number" ? value.frequencySeconds : 0,
    timeoutSeconds: typeof value.timeoutSeconds === "number" ? value.timeoutSeconds : 0,
    status: String(value.status ?? "unmeasured"),
    summary: String(value.summary ?? ""),
    checkCount: typeof value.checkCount === "number" ? value.checkCount : 0,
    createdAt: String(value.createdAt ?? ""),
    updatedAt: String(value.updatedAt ?? ""),
  };
}

function mapCheck(value: Record<string, unknown>): MonitorCheckView {
  return {
    id: String(value.id ?? ""),
    monitorId: String(value.monitorId ?? ""),
    region: String(value.region ?? ""),
    ok: Boolean(value.ok),
    latencyMs: typeof value.latencyMs === "number" ? value.latencyMs : null,
    at: String(value.at ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function mapKey(value: Record<string, unknown>): APIKeyView {
  return {
    id: String(value.id ?? ""),
    name: String(value.name ?? ""),
    prefix: String(value.prefix ?? ""),
    last4: String(value.last4 ?? ""),
    scopes: Array.isArray(value.scopes) ? value.scopes.map((item) => String(item)) : [],
    rateLimitPerMin: typeof value.rateLimitPerMin === "number" ? value.rateLimitPerMin : 0,
    revoked: Boolean(value.revoked),
    createdAt: String(value.createdAt ?? ""),
    lastUsedAt: typeof value.lastUsedAt === "string" ? value.lastUsedAt : null,
  };
}

function mapWebhook(value: Record<string, unknown>): WebhookView {
  return {
    id: String(value.id ?? ""),
    url: String(value.url ?? ""),
    events: Array.isArray(value.events) ? value.events.map((item) => String(item)) : [],
    secretHint: String(value.secretHint ?? ""),
    createdAt: String(value.createdAt ?? ""),
    disabled: Boolean(value.disabled),
  };
}

function mapDelivery(value: Record<string, unknown>): DeliveryView {
  return {
    id: String(value.id ?? ""),
    webhookId: String(value.webhookId ?? ""),
    event: String(value.event ?? ""),
    eventId: String(value.eventId ?? ""),
    timestamp: String(value.timestamp ?? ""),
    idempotencyKey: String(value.idempotencyKey ?? ""),
    signature: String(value.signature ?? ""),
    attempt: typeof value.attempt === "number" ? value.attempt : 0,
    status: String(value.status ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function mapAlert(value: Record<string, unknown>): AlertRuleView {
  return {
    id: String(value.id ?? ""),
    kind: String(value.kind ?? ""),
    monitorId: typeof value.monitorId === "string" ? value.monitorId : null,
    threshold: typeof value.threshold === "number" ? value.threshold : 0,
    enabled: Boolean(value.enabled),
    deliveredCount: typeof value.deliveredCount === "number" ? value.deliveredCount : 0,
    summary: String(value.summary ?? ""),
  };
}

function mapUsage(value: Record<string, unknown>): UsageView {
  return {
    requests: typeof value.requests === "number" ? value.requests : 0,
    measurements: typeof value.measurements === "number" ? value.measurements : 0,
    monitors: typeof value.monitors === "number" ? value.monitors : 0,
    regions: typeof value.regions === "number" ? value.regions : 0,
    webhooks: typeof value.webhooks === "number" ? value.webhooks : 0,
    summary: String(value.summary ?? "Usage counts are stored totals."),
  };
}

function asRecord(value: unknown): Record<string, unknown> {
  return value && typeof value === "object" ? (value as Record<string, unknown>) : {};
}
