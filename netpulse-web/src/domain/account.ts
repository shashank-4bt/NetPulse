export type AccountUser = {
  id: string;
  email: string;
  displayName: string;
  emailVerified: boolean;
  createdAt: string;
  telemetryOptIn: boolean;
};

export type AccountSession = {
  id: string;
  createdAt: string;
  lastSeenAt: string;
  expiresAt: string;
  current: boolean;
  revoked: boolean;
  userAgent: string;
  ip: string;
  label: string;
};

export type SecurityEvent = {
  id: string;
  kind: string;
  at: string;
  summary: string;
  ip: string;
};

export type AuthMethods = {
  password: boolean;
  oauth: string[];
  passkeys: number;
  mfa: string[];
};

export type AuthState = {
  emailSent: boolean;
  emailReason?: string;
  devToken?: string;
  methods: AuthMethods;
};

export type DiagnosisSummary = {
  id: string;
  target: string;
  status: string;
  outcome: string | null;
  createdAt: string;
};

export type UserReport = {
  id: string;
  target: string;
  status: string;
  shared: boolean;
  createdAt: string;
  outcome: string | null;
};

export type SavedService = {
  slug: string;
  createdAt: string;
};

export type AccountDevice = {
  id: string;
  label: string;
  userAgent: string;
  ip: string;
  lastSeenAt: string;
  current: boolean;
  kind: string;
};

export type AlertPreferences = {
  emailEnabled: boolean;
  incidentAlerts: boolean;
  deliveredCount: number;
  summary: string;
};

export type BillingInvoice = {
  id: string;
  amount: string;
  status: string;
};

export type Billing = {
  hasAccount: boolean;
  organizationId: string | null;
  plan: string | null;
  invoices: BillingInvoice[];
  summary: string;
};

export type PrivacySettings = {
  telemetryOptIn: boolean;
  retention: string;
  deletion: string;
  collected: string;
  purpose: string;
  browsingHistory: string;
};

export type ShareLink = {
  token: string;
  path: string;
  summary: string;
};

export type DashboardPayload = {
  internetHealth: string;
  networkInfo: string;
  diagnoses: DiagnosisSummary[];
  savedServices: SavedService[];
  incidents: Array<{
    id: string;
    title: string;
    status: string;
    startedAt: string;
  }>;
  reports: UserReport[];
  alerts: AlertPreferences;
};

export type HistoryQuery = {
  q: string;
  status: string;
  target: string;
  from: string;
  to: string;
};

export const emptyAuthMethods = (): AuthMethods => ({
  password: true,
  oauth: [],
  passkeys: 0,
  mfa: [],
});

export const emptyBilling = (): Billing => ({
  hasAccount: false,
  organizationId: null,
  plan: null,
  invoices: [],
  summary: "No billing account. Organizations and invoices are not enabled.",
});

export const emptyAlerts = (): AlertPreferences => ({
  emailEnabled: false,
  incidentAlerts: false,
  deliveredCount: 0,
  summary: "No alerts have been delivered. Notification delivery is not configured.",
});

export const ACCOUNT_NAV = [
  { href: "/account/profile", label: "Profile" },
  { href: "/account/security", label: "Security" },
  { href: "/account/privacy", label: "Privacy" },
  { href: "/account/devices", label: "Devices" },
  { href: "/account/alerts", label: "Alerts" },
  { href: "/account/billing", label: "Billing" },
] as const;

export const DASHBOARD_NAV = [
  { href: "/dashboard", label: "Overview" },
  { href: "/dashboard/history", label: "History" },
  { href: "/dashboard/reports", label: "Reports" },
  { href: "/developers/dashboard", label: "Developers" },
  { href: "/business/dashboard", label: "Business" },
] as const;
