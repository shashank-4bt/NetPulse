import { emptyObservation, emptyPercentiles, type ObservationView, type PercentileView, type RegionalView } from "@/domain/developer";

export const BUSINESS_NAV = [
  { href: "/business/dashboard", label: "Dashboard" },
  { href: "/business/devices", label: "Devices" },
  { href: "/business/networks", label: "Networks" },
  { href: "/business/services", label: "Services" },
  { href: "/business/incidents", label: "Incidents" },
  { href: "/business/analytics", label: "Analytics" },
  { href: "/business/reports", label: "Reports" },
  { href: "/business/team", label: "Team" },
  { href: "/business/settings", label: "Settings" },
] as const;

export const BUSINESS_ROLES = [
  { value: "owner", label: "Owner" },
  { value: "admin", label: "Admin" },
  { value: "security_admin", label: "Security Admin" },
  { value: "developer", label: "Developer" },
  { value: "analyst", label: "Analyst" },
  { value: "viewer", label: "Viewer" },
  { value: "billing_admin", label: "Billing Admin" },
] as const;

export const REPORT_KINDS = [
  { value: "availability", label: "Availability" },
  { value: "latency", label: "Latency" },
  { value: "incidents", label: "Incidents" },
  { value: "regions", label: "Affected regions" },
  { value: "network", label: "Network performance" },
  { value: "findings", label: "Diagnostic findings" },
] as const;

export const ANALYTICS_FILTERS = ["region", "network", "asn", "service", "endpoint", "device"] as const;

export type OrganizationView = {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
  role: string;
  summary: string;
};

export type MemberView = {
  id: string;
  userId: string;
  email: string;
  displayName: string;
  role: string;
  permissions: string[];
  createdAt: string;
};

export type InviteView = {
  id: string;
  email: string;
  role: string;
  createdAt: string;
  summary: string;
};

export type TeamView = {
  id: string;
  name: string;
  memberIds: string[];
  createdAt: string;
  summary: string;
};

export type OrgDeviceView = {
  id: string;
  name: string;
  label: string;
  region: string;
  networkId: string | null;
  summary: string;
  createdAt: string;
};

export type OrgNetworkView = {
  id: string;
  name: string;
  asn: string;
  region: string;
  summary: string;
  createdAt: string;
};

export type OrgServiceView = {
  id: string;
  name: string;
  slug: string;
  endpoint: string;
  summary: string;
  createdAt: string;
};

export type OrgIncidentView = {
  id: string;
  monitorId: string;
  title: string;
  status: string;
  startedAt: string;
  resolvedAt: string | null;
  deviceIds: string[];
  networkIds: string[];
  serviceIds: string[];
  regions: string[];
  sampleCount: number;
  summary: string;
};

export type OrgDashboardView = {
  overallHealth: string;
  availability: ObservationView;
  incidents: OrgIncidentView[];
  affectedDevices: OrgDeviceView[];
  regions: RegionalView[];
  networks: OrgNetworkView[];
  services: OrgServiceView[];
  summary: string;
};

export type OrgAnalyticsView = {
  filters: Record<string, string>;
  availability: ObservationView;
  latency: PercentileView;
  sampleCount: number;
  incidents: OrgIncidentView[];
  summary: string;
};

export type OrgReportView = {
  id: string;
  kind: string;
  title: string;
  availability: ObservationView;
  latency: PercentileView;
  incidents: OrgIncidentView[];
  regions: RegionalView[];
  networks: OrgNetworkView[];
  findings: string[];
  sampleCount: number;
  createdAt: string;
  summary: string;
};

export type AuditView = {
  id: string;
  actorId: string;
  kind: string;
  at: string;
  summary: string;
};

export type OrgBillingView = {
  hasAccount: boolean;
  organizationId: string | null;
  plan: string | null;
  invoices: unknown[];
  summary: string;
};

export const emptyOrgDashboard = (): OrgDashboardView => ({
  overallHealth: "Not measured. No stored organization checks exist.",
  availability: emptyObservation("Availability"),
  incidents: [],
  affectedDevices: [],
  regions: [],
  networks: [],
  services: [],
  summary: "Organization health is computed from stored checks only. Empty series stay unmeasured.",
});

export const emptyOrgAnalytics = (): OrgAnalyticsView => ({
  filters: {},
  availability: emptyObservation("Availability"),
  latency: emptyPercentiles(),
  sampleCount: 0,
  incidents: [],
  summary: "No stored samples match the selected filters.",
});

export function roleLabel(role: string): string {
  return BUSINESS_ROLES.find((item) => item.value === role)?.label ?? role;
}

export function hasPermission(permissions: string[], need: string): boolean {
  return permissions.includes(need);
}
