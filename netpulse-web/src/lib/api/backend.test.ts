import { afterEach, describe, expect, it, vi } from "vitest";

import { getBackendIncidents, getDiagnosis, isApiConfigured, postDiagnosis } from "@/lib/api/backend";

describe("backend client", () => {
  afterEach(() => {
    vi.unstubAllEnvs();
    vi.unstubAllGlobals();
  });

  it("is unconfigured without NETPULSE_API_BASE_URL", () => {
    vi.stubEnv("NETPULSE_API_BASE_URL", "");
    expect(isApiConfigured()).toBe(false);
  });

  it("maps validation errors without inventing a diagnosis", async () => {
    vi.stubEnv("NETPULSE_API_BASE_URL", "http://api.test");
    vi.stubGlobal(
      "fetch",
      vi.fn(async () =>
        new Response(
          JSON.stringify({
            ok: false,
            error: { code: "ssrf_blocked", message: "blocked host localhost" },
          }),
          { status: 403 }
        )
      )
    );
    const result = await postDiagnosis("localhost");
    expect(result.ok).toBe(false);
    if (!result.ok) {
      expect(result.code).toBe("ssrf_blocked");
    }
  });

  it("maps a completed diagnosis without inventing a cause", async () => {
    vi.stubEnv("NETPULSE_API_BASE_URL", "http://api.test");
    vi.stubGlobal(
      "fetch",
      vi.fn(async () =>
        new Response(
          JSON.stringify({
            ok: true,
            diagnosis: {
              id: "11111111-1111-4111-8111-111111111111",
              status: "insufficient_evidence",
              createdAt: "2026-09-01T00:00:00Z",
              report: {
                reportId: "11111111-1111-4111-8111-111111111111",
                target: {
                  raw: "example.com",
                  hostname: "example.com",
                  kind: "domain",
                  serviceSlug: "",
                },
                timestamp: "2026-09-01T00:00:00Z",
                outcome: "insufficient_evidence",
                tests: [],
                measurements: [],
                evidence: [],
                hypotheses: [],
                alternativeHypotheses: [],
                likelyCause: null,
                confidence: {
                  level: null,
                  percent: null,
                  supportingEvidenceIds: [],
                  alternativeHypothesisIds: [],
                  caveat: "Confidence is not certainty.",
                },
                recommendations: [],
                verificationSteps: [],
                escalationConditions: [],
                graph: [],
                versions: {
                  diagnosticEngineVersion: "0.6.0",
                  ruleVersion: "0.6.0-worker-vantage",
                  measurementVersion: "0.6.0-dns-tcp-tls-http",
                  modelVersion: "0.6.0",
                },
                insufficientEvidence: {
                  determined: true,
                  message: "NetPulse could not safely determine the root cause.",
                  nextCheck: "Compare a user-path measurement.",
                },
                engineVersion: "0.6.0",
              },
            },
          }),
          { status: 200 }
        )
      )
    );
    const result = await getDiagnosis("11111111-1111-4111-8111-111111111111");
    expect(result.ok).toBe(true);
    if (result.ok) {
      expect(result.diagnosis.report?.likelyCause).toBeNull();
      expect(result.diagnosis.report?.target.serviceSlug).toBeNull();
      expect(result.diagnosis.report?.insufficientEvidence.determined).toBe(true);
    }
  });

  it("returns an empty incident list from the API without filling it in", async () => {
    vi.stubEnv("NETPULSE_API_BASE_URL", "http://api.test");
    vi.stubGlobal(
      "fetch",
      vi.fn(async () =>
        new Response(JSON.stringify({ ok: true, incidents: [] }), { status: 200 })
      )
    );
    const result = await getBackendIncidents();
    expect(result.ok).toBe(true);
    if (result.ok) {
      expect(result.incidents).toEqual([]);
    }
  });
});
