import type {
  AccountDevice,
  AccountSession,
  AccountUser,
  AlertPreferences,
  Billing,
  DashboardPayload,
  DiagnosisSummary,
  HistoryQuery,
  PrivacySettings,
  SecurityEvent,
  UserReport,
  SavedService,
} from "@/domain/account";
import {
  apiRequest,
  asApiFailure,
  isApiConfigured,
  type BackendEnvelope,
} from "@/lib/api/backend";
import {
  apiFailure,
  type ApiFailure,
} from "@/lib/api/errors";

export { isApiConfigured };

function isFail(
  value: { status: number; body: BackendEnvelope } | ApiFailure
): value is ApiFailure {
  return "ok" in value && value.ok === false;
}

async function call(
  path: string,
  init: RequestInit & { session?: string | null } = {}
) {
  return apiRequest(path, { timeoutMs: 10000, ...init });
}

export async function getAuthMe(
  session: string
): Promise<{ ok: true; user: AccountUser } | ApiFailure> {
  const result = await call("/v1/auth/me", { session });
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.user) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, user: mapUser(result.body.user) };
}

export async function getDashboard(
  session: string
): Promise<{ ok: true; dashboard: DashboardPayload } | ApiFailure> {
  const result = await call("/v1/me/dashboard", { session });
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.dashboard) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, dashboard: mapDashboard(result.body.dashboard) };
}

export async function getHistory(
  session: string,
  query: HistoryQuery
): Promise<{ ok: true; diagnoses: DiagnosisSummary[] } | ApiFailure> {
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(query)) {
    if (value) {
      params.set(key, value);
    }
  }
  const suffix = params.toString() ? `?${params.toString()}` : "";
  const result = await call(`/v1/me/diagnoses${suffix}`, { session });
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return {
    ok: true,
    diagnoses: (result.body.diagnoses ?? []).map(mapDiagnosisSummary),
  };
}

export async function getReports(
  session: string
): Promise<{ ok: true; reports: UserReport[] } | ApiFailure> {
  const result = await call("/v1/me/reports", { session });
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, reports: (result.body.reports ?? []).map(mapReport) };
}

export async function getSessions(
  session: string
): Promise<{ ok: true; sessions: AccountSession[] } | ApiFailure> {
  const result = await call("/v1/auth/sessions", { session });
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, sessions: (result.body.sessions ?? []).map(mapSession) };
}

export async function getEvents(
  session: string
): Promise<{ ok: true; events: SecurityEvent[] } | ApiFailure> {
  const result = await call("/v1/auth/events", { session });
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, events: (result.body.events ?? []).map(mapEvent) };
}

export async function getPrivacy(
  session: string
): Promise<{ ok: true; privacy: PrivacySettings } | ApiFailure> {
  const result = await call("/v1/me/privacy", { session });
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.privacy) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, privacy: mapPrivacy(result.body.privacy) };
}

export async function getDevices(
  session: string
): Promise<{ ok: true; devices: AccountDevice[] } | ApiFailure> {
  const result = await call("/v1/me/devices", { session });
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, devices: (result.body.devices ?? []).map(mapDevice) };
}

export async function getAlerts(
  session: string
): Promise<{ ok: true; alerts: AlertPreferences } | ApiFailure> {
  const result = await call("/v1/me/alerts", { session });
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.alerts) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, alerts: mapAlerts(result.body.alerts) };
}

export async function getBilling(
  session: string
): Promise<{ ok: true; billing: Billing } | ApiFailure> {
  const result = await call("/v1/me/billing", { session });
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.billing) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, billing: mapBilling(result.body.billing) };
}

export async function getShare(
  token: string
): Promise<{ ok: true; diagnosisId: string } | ApiFailure> {
  const result = await call(`/v1/shares/${encodeURIComponent(token)}`);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.diagnosis?.id) {
    return asApiFailure(result.status || 404, result.body);
  }
  return { ok: true, diagnosisId: String(result.body.diagnosis.id) };
}

export function unavailableAccount(message: string): ApiFailure {
  return apiFailure("unavailable", message, 503);
}

function mapUser(value: Record<string, unknown>): AccountUser {
  return {
    id: String(value.id ?? ""),
    email: String(value.email ?? ""),
    displayName: String(value.displayName ?? ""),
    emailVerified: Boolean(value.emailVerified),
    createdAt: String(value.createdAt ?? ""),
    telemetryOptIn: Boolean(value.telemetryOptIn),
  };
}

function mapDashboard(value: Record<string, unknown>): DashboardPayload {
  return {
    internetHealth: String(value.internetHealth ?? ""),
    networkInfo: String(value.networkInfo ?? ""),
    diagnoses: Array.isArray(value.diagnoses)
      ? value.diagnoses.map((item) => mapDiagnosisSummary(asRecord(item)))
      : [],
    savedServices: Array.isArray(value.savedServices)
      ? value.savedServices.map((item) => mapSaved(asRecord(item)))
      : [],
    incidents: Array.isArray(value.incidents)
      ? value.incidents.map((item) => {
          const record = asRecord(item);
          return {
            id: String(record.id ?? ""),
            title: String(record.title ?? ""),
            status: String(record.status ?? ""),
            startedAt: String(record.startedAt ?? ""),
          };
        })
      : [],
    reports: Array.isArray(value.reports)
      ? value.reports.map((item) => mapReport(asRecord(item)))
      : [],
    alerts: value.alerts ? mapAlerts(asRecord(value.alerts)) : emptyAlertState(),
  };
}

function mapDiagnosisSummary(value: Record<string, unknown>): DiagnosisSummary {
  return {
    id: String(value.id ?? ""),
    target: String(value.target ?? ""),
    status: String(value.status ?? ""),
    outcome: typeof value.outcome === "string" ? value.outcome : null,
    createdAt: String(value.createdAt ?? ""),
  };
}

function mapReport(value: Record<string, unknown>): UserReport {
  return {
    id: String(value.id ?? ""),
    target: String(value.target ?? ""),
    status: String(value.status ?? ""),
    shared: Boolean(value.shared),
    createdAt: String(value.createdAt ?? ""),
    outcome: typeof value.outcome === "string" ? value.outcome : null,
  };
}

function mapSaved(value: Record<string, unknown>): SavedService {
  return {
    slug: String(value.slug ?? ""),
    createdAt: String(value.createdAt ?? ""),
  };
}

function mapSession(value: Record<string, unknown>): AccountSession {
  return {
    id: String(value.id ?? ""),
    createdAt: String(value.createdAt ?? ""),
    lastSeenAt: String(value.lastSeenAt ?? ""),
    expiresAt: String(value.expiresAt ?? ""),
    current: Boolean(value.current),
    revoked: Boolean(value.revoked),
    userAgent: String(value.userAgent ?? ""),
    ip: String(value.ip ?? ""),
    label: String(value.label ?? ""),
  };
}

function mapEvent(value: Record<string, unknown>): SecurityEvent {
  return {
    id: String(value.id ?? ""),
    kind: String(value.kind ?? ""),
    at: String(value.at ?? ""),
    summary: String(value.summary ?? ""),
    ip: String(value.ip ?? ""),
  };
}

function mapPrivacy(value: Record<string, unknown>): PrivacySettings {
  return {
    telemetryOptIn: Boolean(value.telemetryOptIn),
    retention: String(value.retention ?? ""),
    deletion: String(value.deletion ?? ""),
    collected: String(value.collected ?? ""),
    purpose: String(value.purpose ?? ""),
    browsingHistory: String(value.browsingHistory ?? ""),
  };
}

function mapDevice(value: Record<string, unknown>): AccountDevice {
  return {
    id: String(value.id ?? ""),
    label: String(value.label ?? ""),
    userAgent: String(value.userAgent ?? ""),
    ip: String(value.ip ?? ""),
    lastSeenAt: String(value.lastSeenAt ?? value.lastSeen ?? ""),
    current: Boolean(value.current),
    kind: String(value.kind ?? "session"),
  };
}

function mapAlerts(value: Record<string, unknown>): AlertPreferences {
  return {
    emailEnabled: Boolean(value.emailEnabled),
    incidentAlerts: Boolean(value.incidentAlerts),
    deliveredCount:
      typeof value.deliveredCount === "number" ? value.deliveredCount : 0,
    summary: String(value.summary ?? emptyAlertState().summary),
  };
}

function mapBilling(value: Record<string, unknown>): Billing {
  return {
    hasAccount: Boolean(value.hasAccount),
    organizationId:
      typeof value.organizationId === "string" ? value.organizationId : null,
    plan: typeof value.plan === "string" ? value.plan : null,
    invoices: Array.isArray(value.invoices)
      ? value.invoices.map((item) => {
          const record = asRecord(item);
          return {
            id: String(record.id ?? ""),
            amount: String(record.amount ?? ""),
            status: String(record.status ?? ""),
          };
        })
      : [],
    summary: String(value.summary ?? emptyBillingState().summary),
  };
}

function asRecord(value: unknown): Record<string, unknown> {
  return value && typeof value === "object" ? (value as Record<string, unknown>) : {};
}

function emptyAlertState(): AlertPreferences {
  return {
    emailEnabled: false,
    incidentAlerts: false,
    deliveredCount: 0,
    summary: "No alerts have been delivered. Notification delivery is not configured.",
  };
}

function emptyBillingState(): Billing {
  return {
    hasAccount: false,
    organizationId: null,
    plan: null,
    invoices: [],
    summary: "No billing account. Organizations and invoices are not enabled.",
  };
}

