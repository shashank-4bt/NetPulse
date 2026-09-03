import { emptyPercentiles } from "@/domain/developer";
import {
  emptyAdminSystem,
  type AbuseEventView,
  type AdminAuditView,
  type AdminDiagnosisView,
  type AdminIncidentView,
  type AdminMeasurementView,
  type AdminOrgView,
  type AdminServiceView,
  type AdminSystemView,
  type AdminUserView,
  type DiagnosticRuleView,
  type FeatureFlagView,
  type HealthComponentView,
  type OperatorView,
  type RemoteConfigView,
  type RuleOutcomesView,
} from "@/domain/admin";
import { apiRequest, asApiFailure, type BackendEnvelope } from "@/lib/api/backend";
import { type ApiFailure } from "@/lib/api/errors";

function isFail(
  value: { status: number; body: BackendEnvelope } | ApiFailure
): value is ApiFailure {
  return "ok" in value && value.ok === false;
}

async function call(path: string, session: string) {
  return apiRequest(path, { timeoutMs: 10000, session });
}

export async function getAdminMe(
  session: string
): Promise<{ ok: true; operator: OperatorView } | ApiFailure> {
  const result = await call("/v1/admin/me", session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.operator) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, operator: mapOperator(asRecord(result.body.operator)) };
}

export async function getAdminSystem(
  session: string
): Promise<{ ok: true; system: AdminSystemView } | ApiFailure> {
  const result = await call("/v1/admin/system", session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.adminSystem) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, system: mapSystem(asRecord(result.body.adminSystem)) };
}

export async function getAdminUsers(
  session: string
): Promise<{ ok: true; users: AdminUserView[] } | ApiFailure> {
  const result = await list(session, "/v1/admin/users", (body) => (body.adminUsers ?? []).map((item) => mapUser(asRecord(item))));
  return result.ok ? { ok: true, users: result.items } : result;
}

export async function getAdminOrganizations(
  session: string
): Promise<{ ok: true; organizations: AdminOrgView[] } | ApiFailure> {
  const result = await list(session, "/v1/admin/organizations", (body) =>
    (body.organizations ?? []).map((item) => mapOrg(asRecord(item)))
  );
  return result.ok ? { ok: true, organizations: result.items } : result;
}

export async function getAdminServices(
  session: string
): Promise<{ ok: true; services: AdminServiceView[] } | ApiFailure> {
  const result = await list(session, "/v1/admin/services", (body) =>
    (body.services ?? []).map((item) => mapService(asRecord(item)))
  );
  return result.ok ? { ok: true, services: result.items } : result;
}

export async function getAdminIncidents(
  session: string
): Promise<{ ok: true; incidents: AdminIncidentView[] } | ApiFailure> {
  const result = await list(session, "/v1/admin/incidents", (body) =>
    (body.adminIncidents ?? []).map((item) => mapIncident(asRecord(item)))
  );
  return result.ok ? { ok: true, incidents: result.items } : result;
}

export async function getAdminMeasurements(
  session: string
): Promise<{ ok: true; measurements: AdminMeasurementView[] } | ApiFailure> {
  const result = await list(session, "/v1/admin/measurements", (body) =>
    (body.adminMeasurements ?? []).map((item) => mapMeasurement(asRecord(item)))
  );
  return result.ok ? { ok: true, measurements: result.items } : result;
}

export async function getAdminDiagnostics(
  session: string
): Promise<{ ok: true; diagnoses: AdminDiagnosisView[] } | ApiFailure> {
  const result = await list(session, "/v1/admin/diagnostics", (body) =>
    (body.adminDiagnoses ?? []).map((item) => mapDiagnosis(asRecord(item)))
  );
  return result.ok ? { ok: true, diagnoses: result.items } : result;
}

export async function getAdminRules(
  session: string
): Promise<{ ok: true; rules: DiagnosticRuleView[]; outcomes: RuleOutcomesView } | ApiFailure> {
  const result = await call("/v1/admin/rules", session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return {
    ok: true,
    rules: (result.body.adminRules ?? []).map((item) => mapRule(asRecord(item))),
    outcomes: mapOutcomes(asRecord(result.body.ruleOutcomes ?? {})),
  };
}

export async function getAdminAbuse(
  session: string
): Promise<{ ok: true; events: AbuseEventView[] } | ApiFailure> {
  const result = await list(session, "/v1/admin/abuse", (body) =>
    (body.abuseEvents ?? []).map((item) => mapAbuse(asRecord(item)))
  );
  return result.ok ? { ok: true, events: result.items } : result;
}

export async function getAdminAudit(
  session: string
): Promise<{ ok: true; events: AdminAuditView[] } | ApiFailure> {
  const result = await list(session, "/v1/admin/audit", (body) =>
    (body.adminAudit ?? []).map((item) => mapAudit(asRecord(item)))
  );
  return result.ok ? { ok: true, events: result.items } : result;
}

export async function getAdminFlags(
  session: string
): Promise<{ ok: true; flags: FeatureFlagView[] } | ApiFailure> {
  const result = await list(session, "/v1/admin/flags", (body) =>
    (body.featureFlags ?? []).map((item) => mapFlag(asRecord(item)))
  );
  return result.ok ? { ok: true, flags: result.items } : result;
}

export async function getAdminConfig(
  session: string
): Promise<{ ok: true; entries: RemoteConfigView[] } | ApiFailure> {
  const result = await list(session, "/v1/admin/config", (body) =>
    (body.remoteConfig ?? []).map((item) => mapConfig(asRecord(item)))
  );
  return result.ok ? { ok: true, entries: result.items } : result;
}

async function list<T>(
  session: string,
  path: string,
  map: (body: BackendEnvelope) => T[]
): Promise<{ ok: true; items: T[] } | ApiFailure> {
  const result = await call(path, session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, items: map(result.body) };
}

export function mapSystem(value: Record<string, unknown>): AdminSystemView {
  const empty = emptyAdminSystem();
  return {
    api: mapHealth(asRecord(value.api), empty.api),
    worker: mapHealth(asRecord(value.worker), empty.worker),
    queue: mapHealth(asRecord(value.queue), empty.queue),
    database: mapHealth(asRecord(value.database), empty.database),
    cache: mapHealth(asRecord(value.cache), empty.cache),
    measurementFailures: mapObservation(asRecord(value.measurementFailures), empty.measurementFailures),
    errorRates: mapObservation(asRecord(value.errorRates), empty.errorRates),
    latency: mapPercentiles(asRecord(value.latency)),
    summary: String(value.summary ?? empty.summary),
  };
}

function mapHealth(value: Record<string, unknown>, fallback: HealthComponentView): HealthComponentView {
  return {
    name: String(value.name ?? fallback.name),
    status: String(value.status ?? fallback.status),
    detail: String(value.detail ?? fallback.detail),
    measured: Boolean(value.measured),
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

function mapOperator(value: Record<string, unknown>): OperatorView {
  return {
    userId: String(value.userId ?? ""),
    email: String(value.email ?? ""),
    role: String(value.role ?? "operator"),
    permissions: Array.isArray(value.permissions) ? value.permissions.map((item) => String(item)) : [],
  };
}

function mapUser(value: Record<string, unknown>): AdminUserView {
  return {
    id: String(value.id ?? ""),
    email: String(value.email ?? ""),
    displayName: String(value.displayName ?? ""),
    emailVerified: Boolean(value.emailVerified),
    createdAt: String(value.createdAt ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function mapOrg(value: Record<string, unknown>): AdminOrgView {
  return {
    id: String(value.id ?? ""),
    name: String(value.name ?? ""),
    createdAt: String(value.createdAt ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function mapService(value: Record<string, unknown>): AdminServiceView {
  return {
    slug: String(value.slug ?? ""),
    name: String(value.name ?? ""),
    category: String(value.category ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function mapIncident(value: Record<string, unknown>): AdminIncidentView {
  const override =
    value.override && typeof value.override === "object" && !Array.isArray(value.override)
      ? (value.override as Record<string, unknown>)
      : null;
  return {
    id: String(value.id ?? ""),
    title: String(value.title ?? ""),
    severity: String(value.severity ?? ""),
    status: String(value.status ?? ""),
    startedAt: String(value.startedAt ?? ""),
    lastUpdatedAt: String(value.lastUpdatedAt ?? ""),
    sampleCount: typeof value.sampleCount === "number" ? value.sampleCount : 0,
    notes: Array.isArray(value.notes)
      ? value.notes.map((item) => {
          const note = asRecord(item);
          return {
            id: String(note.id ?? ""),
            incidentId: String(note.incidentId ?? ""),
            kind: String(note.kind ?? ""),
            body: String(note.body ?? ""),
            actorId: String(note.actorId ?? ""),
            at: String(note.at ?? ""),
          };
        })
      : [],
    overrideClassification: override ? String(override.classification ?? "") : null,
    overrideReason: override ? String(override.reason ?? "") : null,
  };
}

function mapMeasurement(value: Record<string, unknown>): AdminMeasurementView {
  return {
    diagnosisId: String(value.diagnosisId ?? ""),
    key: String(value.key ?? ""),
    label: String(value.label ?? ""),
    measured: Boolean(value.measured),
    summary: String(value.summary ?? ""),
  };
}

function mapDiagnosis(value: Record<string, unknown>): AdminDiagnosisView {
  return {
    id: String(value.id ?? ""),
    target: String(value.target ?? ""),
    status: String(value.status ?? ""),
    userId: String(value.userId ?? ""),
    createdAt: String(value.createdAt ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function mapRule(value: Record<string, unknown>): DiagnosticRuleView {
  const thresholds =
    value.thresholds && typeof value.thresholds === "object" && !Array.isArray(value.thresholds)
      ? Object.fromEntries(
          Object.entries(value.thresholds as Record<string, unknown>).map(([key, item]) => [key, String(item)])
        )
      : {};
  return {
    id: String(value.id ?? ""),
    name: String(value.name ?? ""),
    version: String(value.version ?? ""),
    layer: String(value.layer ?? ""),
    thresholds,
    summary: String(value.summary ?? ""),
  };
}

function mapOutcomes(value: Record<string, unknown>): RuleOutcomesView {
  const statuses =
    value.statuses && typeof value.statuses === "object" && !Array.isArray(value.statuses)
      ? Object.fromEntries(
          Object.entries(value.statuses as Record<string, unknown>).map(([key, item]) => [
            key,
            typeof item === "number" ? item : 0,
          ])
        )
      : {};
  return {
    statuses,
    falsePositives: typeof value.falsePositives === "number" ? value.falsePositives : 0,
    falseNegatives: typeof value.falseNegatives === "number" ? value.falseNegatives : 0,
    sampleCount: typeof value.sampleCount === "number" ? value.sampleCount : 0,
    summary: String(value.summary ?? "False positives and false negatives stay 0 until an operator stores a label."),
  };
}

function mapAbuse(value: Record<string, unknown>): AbuseEventView {
  return {
    id: String(value.id ?? ""),
    kind: String(value.kind ?? ""),
    actor: String(value.actor ?? ""),
    ip: String(value.ip ?? ""),
    resource: String(value.resource ?? ""),
    at: String(value.at ?? ""),
    result: String(value.result ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function mapAudit(value: Record<string, unknown>): AdminAuditView {
  return {
    id: String(value.id ?? ""),
    actorId: String(value.actorId ?? ""),
    action: String(value.action ?? ""),
    resource: String(value.resource ?? ""),
    at: String(value.at ?? ""),
    result: String(value.result ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function mapFlag(value: Record<string, unknown>): FeatureFlagView {
  return {
    id: String(value.id ?? ""),
    name: String(value.name ?? ""),
    environment: String(value.environment ?? ""),
    enabled: Boolean(value.enabled),
    percentage: typeof value.percentage === "number" ? value.percentage : 0,
    userIds: Array.isArray(value.userIds) ? value.userIds.map((item) => String(item)) : [],
    orgIds: Array.isArray(value.orgIds) ? value.orgIds.map((item) => String(item)) : [],
    updatedAt: String(value.updatedAt ?? ""),
    summary: String(value.summary ?? ""),
    targetMatch: typeof value.targetMatch === "boolean" ? value.targetMatch : null,
  };
}

function mapConfig(value: Record<string, unknown>): RemoteConfigView {
  return {
    key: String(value.key ?? ""),
    value: String(value.value ?? ""),
    updatedAt: String(value.updatedAt ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function asRecord(value: unknown): Record<string, unknown> {
  return value && typeof value === "object" && !Array.isArray(value) ? (value as Record<string, unknown>) : {};
}
