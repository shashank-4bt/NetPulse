import type { DiagnosticReport } from "@/domain/diagnostic";
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

export type BackendIncident = {
  id: string;
  title: string;
  scope: string;
  startedAt: string;
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
  incidents?: BackendIncident[];
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

export async function getBackendIncidents(): Promise<
  { ok: true; incidents: BackendIncident[] } | ApiFailure
> {
  const result = await request("/v1/incidents");
  if (isApiFailure(result)) {
    return result;
  }
  const { status, body } = result;
  if (status >= 400) {
    return asFailure(status, body);
  }
  return { ok: true, incidents: body.incidents ?? [] };
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
