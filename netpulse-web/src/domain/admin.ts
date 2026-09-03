import { emptyObservation, emptyPercentiles, type ObservationView, type PercentileView } from "@/domain/developer";

export const ADMIN_NAV = [
  { href: "/admin", label: "Overview" },
  { href: "/admin/users", label: "Users" },
  { href: "/admin/organizations", label: "Organizations" },
  { href: "/admin/services", label: "Services" },
  { href: "/admin/incidents", label: "Incidents" },
  { href: "/admin/measurements", label: "Measurements" },
  { href: "/admin/diagnostics", label: "Diagnostics" },
  { href: "/admin/rules", label: "Rules" },
  { href: "/admin/abuse", label: "Abuse" },
  { href: "/admin/audit", label: "Audit" },
  { href: "/admin/system", label: "System" },
] as const;

export type HealthComponentView = {
  name: string;
  status: string;
  detail: string;
  measured: boolean;
};

export type AdminSystemView = {
  api: HealthComponentView;
  worker: HealthComponentView;
  queue: HealthComponentView;
  database: HealthComponentView;
  cache: HealthComponentView;
  measurementFailures: ObservationView;
  errorRates: ObservationView;
  latency: PercentileView;
  summary: string;
};

export type OperatorView = {
  userId: string;
  email: string;
  role: string;
  permissions: string[];
};

export type AdminUserView = {
  id: string;
  email: string;
  displayName: string;
  emailVerified: boolean;
  createdAt: string;
  summary: string;
};

export type AdminMeasurementView = {
  diagnosisId: string;
  key: string;
  label: string;
  measured: boolean;
  summary: string;
};

export type AdminDiagnosisView = {
  id: string;
  target: string;
  status: string;
  userId: string;
  createdAt: string;
  summary: string;
};

export type DiagnosticRuleView = {
  id: string;
  name: string;
  version: string;
  layer: string;
  thresholds: Record<string, string>;
  summary: string;
};

export type RuleOutcomesView = {
  statuses: Record<string, number>;
  falsePositives: number;
  falseNegatives: number;
  sampleCount: number;
  summary: string;
};

export type AbuseEventView = {
  id: string;
  kind: string;
  actor: string;
  ip: string;
  resource: string;
  at: string;
  result: string;
  summary: string;
};

export type AdminAuditView = {
  id: string;
  actorId: string;
  action: string;
  resource: string;
  at: string;
  result: string;
  summary: string;
};

export type FeatureFlagView = {
  id: string;
  name: string;
  environment: string;
  enabled: boolean;
  percentage: number;
  userIds: string[];
  orgIds: string[];
  updatedAt: string;
  summary: string;
  targetMatch: boolean | null;
};

export type RemoteConfigView = {
  key: string;
  value: string;
  updatedAt: string;
  summary: string;
};

export type IncidentNoteView = {
  id: string;
  incidentId: string;
  kind: string;
  body: string;
  actorId: string;
  at: string;
};

export type AdminIncidentView = {
  id: string;
  title: string;
  severity: string;
  status: string;
  startedAt: string;
  lastUpdatedAt: string;
  sampleCount: number;
  notes: IncidentNoteView[];
  overrideClassification: string | null;
  overrideReason: string | null;
};

export type AdminServiceView = {
  slug: string;
  name: string;
  category: string;
  summary: string;
};

export type AdminOrgView = {
  id: string;
  name: string;
  createdAt: string;
  summary: string;
};

export function emptyHealth(name: string, detail: string): HealthComponentView {
  return { name, status: "unmeasured", detail, measured: false };
}

export function emptyAdminSystem(): AdminSystemView {
  return {
    api: { name: "API", status: "up", detail: "This process answered the request.", measured: true },
    worker: emptyHealth("Worker", "No worker heartbeat is stored."),
    queue: { name: "Queue", status: "observed", detail: "Queue depth: 0.", measured: true },
    database: emptyHealth("Database", "No live database probe is stored."),
    cache: emptyHealth("Cache", "No live cache probe is stored."),
    measurementFailures: emptyObservation("Measurement failures"),
    errorRates: emptyObservation("Error rate"),
    latency: emptyPercentiles(),
    summary: "System figures use stored process state and stored measurements only. Empty series stay unmeasured.",
  };
}
