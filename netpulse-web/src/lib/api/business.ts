import type { APIKeyView, MonitorView } from "@/domain/developer";
import { emptyObservation, emptyPercentiles } from "@/domain/developer";
import {
  emptyOrgAnalytics,
  emptyOrgDashboard,
  type AuditView,
  type InviteView,
  type MemberView,
  type OrgAnalyticsView,
  type OrgBillingView,
  type OrgDashboardView,
  type OrgDeviceView,
  type OrgIncidentView,
  type OrgNetworkView,
  type OrganizationView,
  type OrgReportView,
  type OrgServiceView,
  type TeamView,
} from "@/domain/business";
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

export async function listOrgs(
  session: string
): Promise<{ ok: true; organizations: OrganizationView[] } | ApiFailure> {
  const result = await call("/v1/orgs", session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, organizations: (result.body.organizations ?? []).map((item) => mapOrg(asRecord(item))) };
}

export async function getOrg(
  session: string,
  orgId: string
): Promise<{ ok: true; organization: OrganizationView; permissions: string[] } | ApiFailure> {
  const result = await call(`/v1/orgs/${orgId}`, session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.organization) {
    return asApiFailure(result.status, result.body);
  }
  return {
    ok: true,
    organization: mapOrg(asRecord(result.body.organization)),
    permissions: result.body.permissions ?? [],
  };
}

export async function getOrgDashboard(
  session: string,
  orgId: string
): Promise<{ ok: true; dashboard: OrgDashboardView } | ApiFailure> {
  const result = await call(`/v1/orgs/${orgId}/dashboard`, session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.orgDashboard) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, dashboard: mapDashboard(asRecord(result.body.orgDashboard)) };
}

export async function getOrgAnalytics(
  session: string,
  orgId: string,
  filters: Record<string, string>
): Promise<{ ok: true; analytics: OrgAnalyticsView } | ApiFailure> {
  const query = new URLSearchParams();
  for (const [key, value] of Object.entries(filters)) {
    if (value) {
      query.set(key, value);
    }
  }
  const suffix = query.size ? `?${query.toString()}` : "";
  const result = await call(`/v1/orgs/${orgId}/analytics${suffix}`, session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.analytics) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, analytics: mapAnalytics(asRecord(result.body.analytics)) };
}

export async function getOrgMembers(
  session: string,
  orgId: string
): Promise<{ ok: true; members: MemberView[] } | ApiFailure> {
  const items = await listItems(session, `/v1/orgs/${orgId}/members`, (body) =>
    (body.members ?? []).map((item) => mapMember(asRecord(item)))
  );
  return items.ok ? { ok: true, members: items.items } : items;
}

export async function getOrgInvites(
  session: string,
  orgId: string
): Promise<{ ok: true; invites: InviteView[] } | ApiFailure> {
  const items = await listItems(session, `/v1/orgs/${orgId}/invites`, (body) =>
    (body.invites ?? []).map((item) => mapInvite(asRecord(item)))
  );
  return items.ok ? { ok: true, invites: items.items } : items;
}

export async function getOrgTeams(
  session: string,
  orgId: string
): Promise<{ ok: true; teams: TeamView[] } | ApiFailure> {
  const items = await listItems(session, `/v1/orgs/${orgId}/teams`, (body) =>
    (body.teams ?? []).map((item) => mapTeam(asRecord(item)))
  );
  return items.ok ? { ok: true, teams: items.items } : items;
}

export async function getOrgDevices(
  session: string,
  orgId: string
): Promise<{ ok: true; devices: OrgDeviceView[] } | ApiFailure> {
  const items = await listItems(session, `/v1/orgs/${orgId}/devices`, (body) =>
    (body.orgDevices ?? []).map((item) => mapDevice(asRecord(item)))
  );
  return items.ok ? { ok: true, devices: items.items } : items;
}

export async function getOrgNetworks(
  session: string,
  orgId: string
): Promise<{ ok: true; networks: OrgNetworkView[] } | ApiFailure> {
  const items = await listItems(session, `/v1/orgs/${orgId}/networks`, (body) =>
    (body.networks ?? []).map((item) => mapNetwork(asRecord(item)))
  );
  return items.ok ? { ok: true, networks: items.items } : items;
}

export async function getOrgServices(
  session: string,
  orgId: string
): Promise<{ ok: true; services: OrgServiceView[] } | ApiFailure> {
  const items = await listItems(session, `/v1/orgs/${orgId}/services`, (body) =>
    (body.orgServices ?? []).map((item) => mapService(asRecord(item)))
  );
  return items.ok ? { ok: true, services: items.items } : items;
}

export async function getOrgMonitors(
  session: string,
  orgId: string
): Promise<{ ok: true; monitors: MonitorView[] } | ApiFailure> {
  const items = await listItems(session, `/v1/orgs/${orgId}/monitors`, (body) =>
    (body.monitors ?? []).map((item) => mapMonitor(asRecord(item)))
  );
  return items.ok ? { ok: true, monitors: items.items } : items;
}

export async function getOrgIncidents(
  session: string,
  orgId: string
): Promise<{ ok: true; incidents: OrgIncidentView[] } | ApiFailure> {
  const items = await listItems(session, `/v1/orgs/${orgId}/incidents`, (body) =>
    (body.orgIncidents ?? []).map((item) => mapIncident(asRecord(item)))
  );
  return items.ok ? { ok: true, incidents: items.items } : items;
}

export async function getOrgReports(
  session: string,
  orgId: string
): Promise<{ ok: true; reports: OrgReportView[] } | ApiFailure> {
  const items = await listItems(session, `/v1/orgs/${orgId}/reports`, (body) =>
    (body.orgReports ?? []).map((item) => mapReport(asRecord(item)))
  );
  return items.ok ? { ok: true, reports: items.items } : items;
}

export async function getOrgKeys(
  session: string,
  orgId: string
): Promise<{ ok: true; keys: APIKeyView[] } | ApiFailure> {
  const items = await listItems(session, `/v1/orgs/${orgId}/keys`, (body) =>
    (body.apiKeys ?? []).map((item) => mapKey(asRecord(item)))
  );
  return items.ok ? { ok: true, keys: items.items } : items;
}

export async function getOrgBilling(
  session: string,
  orgId: string
): Promise<{ ok: true; billing: OrgBillingView } | ApiFailure> {
  const result = await call(`/v1/orgs/${orgId}/billing`, session);
  if (isFail(result)) {
    return result;
  }
  if (result.status >= 400 || !result.body.billing) {
    return asApiFailure(result.status, result.body);
  }
  return { ok: true, billing: mapBilling(asRecord(result.body.billing)) };
}

export async function getOrgAudit(
  session: string,
  orgId: string
): Promise<{ ok: true; events: AuditView[] } | ApiFailure> {
  const items = await listItems(session, `/v1/orgs/${orgId}/audit`, (body) =>
    (body.auditEvents ?? []).map((item) => mapAudit(asRecord(item)))
  );
  return items.ok ? { ok: true, events: items.items } : items;
}

async function listItems<T>(
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

export function mapDashboard(value: Record<string, unknown>): OrgDashboardView {
  const empty = emptyOrgDashboard();
  return {
    overallHealth: String(value.overallHealth ?? empty.overallHealth),
    availability: mapObservation(asRecord(value.availability), empty.availability),
    incidents: Array.isArray(value.incidents) ? value.incidents.map((item) => mapIncident(asRecord(item))) : [],
    affectedDevices: Array.isArray(value.affectedDevices)
      ? value.affectedDevices.map((item) => mapDevice(asRecord(item)))
      : [],
    regions: Array.isArray(value.regions) ? value.regions.map((item) => mapRegional(asRecord(item))) : [],
    networks: Array.isArray(value.networks) ? value.networks.map((item) => mapNetwork(asRecord(item))) : [],
    services: Array.isArray(value.services) ? value.services.map((item) => mapService(asRecord(item))) : [],
    summary: String(value.summary ?? empty.summary),
  };
}

export function mapAnalytics(value: Record<string, unknown>): OrgAnalyticsView {
  const empty = emptyOrgAnalytics();
  const filters =
    value.filters && typeof value.filters === "object" && !Array.isArray(value.filters)
      ? Object.fromEntries(Object.entries(value.filters as Record<string, unknown>).map(([key, item]) => [key, String(item)]))
      : {};
  return {
    filters,
    availability: mapObservation(asRecord(value.availability), empty.availability),
    latency: mapPercentiles(asRecord(value.latency)),
    sampleCount: typeof value.sampleCount === "number" ? value.sampleCount : 0,
    incidents: Array.isArray(value.incidents) ? value.incidents.map((item) => mapIncident(asRecord(item))) : [],
    summary: String(value.summary ?? empty.summary),
  };
}

export function mapReport(value: Record<string, unknown>): OrgReportView {
  return {
    id: String(value.id ?? ""),
    kind: String(value.kind ?? ""),
    title: String(value.title ?? ""),
    availability: mapObservation(asRecord(value.availability), emptyObservation("Availability")),
    latency: mapPercentiles(asRecord(value.latency)),
    incidents: Array.isArray(value.incidents) ? value.incidents.map((item) => mapIncident(asRecord(item))) : [],
    regions: Array.isArray(value.regions) ? value.regions.map((item) => mapRegional(asRecord(item))) : [],
    networks: Array.isArray(value.networks) ? value.networks.map((item) => mapNetwork(asRecord(item))) : [],
    findings: Array.isArray(value.findings) ? value.findings.map((item) => String(item)) : [],
    sampleCount: typeof value.sampleCount === "number" ? value.sampleCount : 0,
    createdAt: String(value.createdAt ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function mapOrg(value: Record<string, unknown>): OrganizationView {
  return {
    id: String(value.id ?? ""),
    name: String(value.name ?? ""),
    createdAt: String(value.createdAt ?? ""),
    updatedAt: String(value.updatedAt ?? ""),
    role: String(value.role ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function mapMember(value: Record<string, unknown>): MemberView {
  return {
    id: String(value.id ?? ""),
    userId: String(value.userId ?? ""),
    email: String(value.email ?? ""),
    displayName: String(value.displayName ?? ""),
    role: String(value.role ?? ""),
    permissions: Array.isArray(value.permissions) ? value.permissions.map((item) => String(item)) : [],
    createdAt: String(value.createdAt ?? ""),
  };
}

function mapInvite(value: Record<string, unknown>): InviteView {
  return {
    id: String(value.id ?? ""),
    email: String(value.email ?? ""),
    role: String(value.role ?? ""),
    createdAt: String(value.createdAt ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function mapTeam(value: Record<string, unknown>): TeamView {
  return {
    id: String(value.id ?? ""),
    name: String(value.name ?? ""),
    memberIds: Array.isArray(value.memberIds) ? value.memberIds.map((item) => String(item)) : [],
    createdAt: String(value.createdAt ?? ""),
    summary: String(value.summary ?? ""),
  };
}

function mapDevice(value: Record<string, unknown>): OrgDeviceView {
  return {
    id: String(value.id ?? ""),
    name: String(value.name ?? ""),
    label: String(value.label ?? ""),
    region: String(value.region ?? ""),
    networkId: typeof value.networkId === "string" ? value.networkId : null,
    summary: String(value.summary ?? ""),
    createdAt: String(value.createdAt ?? ""),
  };
}

function mapNetwork(value: Record<string, unknown>): OrgNetworkView {
  return {
    id: String(value.id ?? ""),
    name: String(value.name ?? ""),
    asn: String(value.asn ?? ""),
    region: String(value.region ?? ""),
    summary: String(value.summary ?? ""),
    createdAt: String(value.createdAt ?? ""),
  };
}

function mapService(value: Record<string, unknown>): OrgServiceView {
  return {
    id: String(value.id ?? ""),
    name: String(value.name ?? ""),
    slug: String(value.slug ?? ""),
    endpoint: String(value.endpoint ?? ""),
    summary: String(value.summary ?? ""),
    createdAt: String(value.createdAt ?? ""),
  };
}

function mapIncident(value: Record<string, unknown>): OrgIncidentView {
  return {
    id: String(value.id ?? ""),
    monitorId: String(value.monitorId ?? ""),
    title: String(value.title ?? ""),
    status: String(value.status ?? ""),
    startedAt: String(value.startedAt ?? ""),
    resolvedAt: typeof value.resolvedAt === "string" ? value.resolvedAt : null,
    deviceIds: Array.isArray(value.deviceIds) ? value.deviceIds.map((item) => String(item)) : [],
    networkIds: Array.isArray(value.networkIds) ? value.networkIds.map((item) => String(item)) : [],
    serviceIds: Array.isArray(value.serviceIds) ? value.serviceIds.map((item) => String(item)) : [],
    regions: Array.isArray(value.regions) ? value.regions.map((item) => String(item)) : [],
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

function mapBilling(value: Record<string, unknown>): OrgBillingView {
  return {
    hasAccount: Boolean(value.hasAccount),
    organizationId: typeof value.organizationId === "string" ? value.organizationId : null,
    plan: typeof value.plan === "string" ? value.plan : null,
    invoices: Array.isArray(value.invoices) ? value.invoices : [],
    summary: String(value.summary ?? "No billing account is stored for this organization."),
  };
}

function mapAudit(value: Record<string, unknown>): AuditView {
  return {
    id: String(value.id ?? ""),
    actorId: String(value.actorId ?? ""),
    kind: String(value.kind ?? ""),
    at: String(value.at ?? ""),
    summary: String(value.summary ?? ""),
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

function asRecord(value: unknown): Record<string, unknown> {
  return value && typeof value === "object" ? (value as Record<string, unknown>) : {};
}
