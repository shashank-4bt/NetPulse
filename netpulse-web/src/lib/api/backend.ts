import type { DiagnosticReport } from "@/domain/diagnostic";
import type {
  PublicIncidentRecord,
  ServiceIntelligence,
} from "@/domain/observatory";
import { emptyIncidentLists, emptyServiceIntelligence } from "@/features/observatory/empty-intelligence";
import {
  apiFailure,
  classifyTransportError,
  parseErrorCode,
  type ApiFailure,
} from "@/lib/api/errors";

export type DiagnosisStatus =
  | "queued"
  | "running"
  | "complete"
  | "partial"
  | "failed"
  | "unavailable"
  | "insufficient_evidence";

export type BackendDiagnosis = {
  id: string;
  status: DiagnosisStatus;
  createdAt: string;
  report: DiagnosticReport | null;
};

export type BackendService = {
  slug: string;
  name: string;
  category: string;
  summary: string;
  layers: string[];
};

export type BackendIncident = PublicIncidentRecord;

export type IncidentListQuery = {
  service?: string;
  region?: string;
  network?: string;
  severity?: string;
  status?: string;
  time?: string;
  q?: string;
  sort?: string;
  page?: number;
  pageSize?: number;
};

export function getApiBaseUrl(): string | null {
  const value = process.env.NETPULSE_API_BASE_URL?.trim();
  return value ? value.replace(/\/$/, "") : null;
}

export function isApiConfigured(): boolean {
  return getApiBaseUrl() !== null;
}

type Envelope = {
  ok?: boolean;
  diagnosis?: {
    id?: string;
    status?: string;
    createdAt?: string;
    report?: DiagnosticReport | null;
  };
  services?: BackendService[];
  service?: BackendService;
  intelligence?: ServiceIntelligence;
  incidents?: Record<string, unknown>[];
  incident?: Record<string, unknown>;
  page?: { number?: number; size?: number; total?: number };
  health?: {
    status?: string;
    version?: string;
    storage?: Record<string, string>;
  };
  error?: { code?: string; message?: string };
};

function isApiFailure(
  value: { status: number; body: Envelope } | ApiFailure
): value is ApiFailure {
  return "ok" in value && value.ok === false;
}

async function request(
  path: string,
  init: RequestInit & { timeoutMs?: number } = {}
): Promise<{ status: number; body: Envelope } | ApiFailure> {
  const base = getApiBaseUrl();
  if (!base) {
    return apiFailure("unavailable", "NETPULSE_API_BASE_URL is not set.", 503);
  }
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), init.timeoutMs ?? 10000);
  try {
    const response = await fetch(`${base}${path}`, {
      ...init,
      signal: controller.signal,
      headers: {
        Accept: "application/json",
        ...(init.body ? { "Content-Type": "application/json" } : {}),
        ...init.headers,
      },
      cache: "no-store",
    });
    let body: Envelope = {};
    try {
      body = (await response.json()) as Envelope;
    } catch {
      body = {};
    }
    return { status: response.status, body };
  } catch (error) {
    return classifyTransportError(error);
  } finally {
    clearTimeout(timeout);
  }
}

function asFailure(status: number, body: Envelope): ApiFailure {
  return apiFailure(
    parseErrorCode(body.error?.code),
    body.error?.message ?? "Request failed.",
    status
  );
}

export async function postDiagnosis(
  target: string
): Promise<{ ok: true; diagnosis: BackendDiagnosis } | ApiFailure> {
  const result = await request("/v1/diagnoses", {
    method: "POST",
    body: JSON.stringify({ target }),
    timeoutMs: 12000,
  });
  if (isApiFailure(result)) {
    return result;
  }
  const { status, body } = result;
  if (status >= 400 || !body.diagnosis?.id) {
    return asFailure(status, body);
  }
  return { ok: true, diagnosis: mapDiagnosis(body.diagnosis) };
}

export async function getDiagnosis(
  id: string
): Promise<{ ok: true; diagnosis: BackendDiagnosis } | ApiFailure> {
  const result = await request(`/v1/diagnoses/${id}`, { timeoutMs: 8000 });
  if (isApiFailure(result)) {
    return result;
  }
  const { status, body } = result;
  if (status >= 400 || !body.diagnosis?.id) {
    return asFailure(status, body);
  }
  return { ok: true, diagnosis: mapDiagnosis(body.diagnosis) };
}

export async function getBackendIncidents(
  query: IncidentListQuery = {}
): Promise<
  | { ok: true; incidents: BackendIncident[]; total: number }
  | ApiFailure
> {
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(query)) {
    if (value === undefined || value === "") {
      continue;
    }
    params.set(key, String(value));
  }
  const encoded = params.toString();
  const suffix = encoded ? `?${encoded}` : "";
  const result = await request(`/v1/incidents${suffix}`);
  if (isApiFailure(result)) {
    return result;
  }
  const { status, body } = result;
  if (status >= 400) {
    return asFailure(status, body);
  }
  const incidents = (body.incidents ?? []).map(mapIncident);
  return {
    ok: true,
    incidents,
    total: body.page?.total ?? incidents.length,
  };
}

export async function getBackendIncident(
  id: string
): Promise<{ ok: true; incident: BackendIncident } | ApiFailure> {
  const result = await request(`/v1/incidents/${id}`);
  if (isApiFailure(result)) {
    return result;
  }
  const { status, body } = result;
  if (status === 404) {
    return asFailure(status, body);
  }
  if (status >= 400 || !body.incident) {
    return asFailure(status, body);
  }
  return { ok: true, incident: mapIncident(body.incident) };
}

export async function getBackendService(
  slug: string
): Promise<
  | { ok: true; service: BackendService; intelligence: ServiceIntelligence }
  | ApiFailure
> {
  const result = await request(`/v1/services/${slug}`);
  if (isApiFailure(result)) {
    return result;
  }
  const { status, body } = result;
  if (status >= 400 || !body.service?.slug) {
    return asFailure(status, body);
  }
  return {
    ok: true,
    service: body.service,
    intelligence: mapIntelligence(body.intelligence),
  };
}

export async function getBackendServices(): Promise<
  { ok: true; services: BackendService[] } | ApiFailure
> {
  const result = await request("/v1/services");
  if (isApiFailure(result)) {
    return result;
  }
  const { status, body } = result;
  if (status >= 400) {
    return asFailure(status, body);
  }
  return { ok: true, services: body.services ?? [] };
}

function mapIntelligence(value: Envelope["intelligence"]): ServiceIntelligence {
  const empty = emptyServiceIntelligence();
  if (!value) {
    return empty;
  }
  return {
    currentState: value.currentState ?? empty.currentState,
    health: value.health ?? null,
    lastUpdated: value.lastUpdated ?? null,
    availability: {
      ...empty.availability,
      ...value.availability,
      value: value.availability?.value ?? null,
      measured: Boolean(value.availability?.measured),
      sampleCount: value.availability?.sampleCount ?? 0,
    },
    latency: {
      ...empty.latency,
      ...value.latency,
      value: value.latency?.value ?? null,
      measured: Boolean(value.latency?.measured),
      sampleCount: value.latency?.sampleCount ?? 0,
    },
    errors: {
      ...empty.errors,
      ...value.errors,
      value: value.errors?.value ?? null,
      measured: Boolean(value.errors?.measured),
      sampleCount: value.errors?.sampleCount ?? 0,
    },
    regionalHealth: value.regionalHealth ?? [],
    networkHealth: value.networkHealth ?? [],
    recentIncidentIds: value.recentIncidentIds ?? [],
  };
}

function asStringList(value: unknown): string[] {
  if (!Array.isArray(value)) {
    return [];
  }
  return value.filter((item): item is string => typeof item === "string");
}

function mapIncident(value: Record<string, unknown>): PublicIncidentRecord {
  const lists = emptyIncidentLists();
  return {
    id: String(value.id ?? ""),
    title: String(value.title ?? ""),
    severity: (typeof value.severity === "string" ? value.severity : "") as PublicIncidentRecord["severity"],
    status: (typeof value.status === "string" ? value.status : "") as PublicIncidentRecord["status"],
    scope: String(value.scope ?? ""),
    startedAt: String(value.startedAt ?? ""),
    lastUpdatedAt: String(value.lastUpdatedAt ?? value.startedAt ?? ""),
    affectedServices: asStringList(value.affectedServices),
    regions: asStringList(value.regions),
    networks: asStringList(value.networks),
    evidence: Array.isArray(value.evidence) ? (value.evidence as PublicIncidentRecord["evidence"]) : lists.evidence,
    hypotheses: Array.isArray(value.hypotheses)
      ? (value.hypotheses as PublicIncidentRecord["hypotheses"])
      : lists.hypotheses,
    confidence:
      value.confidence && typeof value.confidence === "object"
        ? {
            ...lists.confidence,
            ...(value.confidence as Partial<PublicIncidentRecord["confidence"]>),
          }
        : lists.confidence,
    timeline: Array.isArray(value.timeline)
      ? (value.timeline as PublicIncidentRecord["timeline"])
      : lists.timeline,
    sampleCount: typeof value.sampleCount === "number" ? value.sampleCount : 0,
    sampleRate: typeof value.sampleRate === "string" ? value.sampleRate : null,
    affectedUserCount: null,
  };
}

function mapDiagnosis(value: NonNullable<Envelope["diagnosis"]>): BackendDiagnosis {
  return {
    id: value.id ?? "",
    status: (value.status ?? "unknown") as DiagnosisStatus,
    createdAt: value.createdAt ?? "",
    report: value.report
      ? {
          ...value.report,
          target: {
            ...value.report.target,
            serviceSlug: value.report.target.serviceSlug || null,
          },
        }
      : null,
  };
}
