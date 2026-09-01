import { describe, expect, it } from "vitest";

import { GET } from "@/app/api/reports/[id]/route";
import { POST } from "@/app/api/reports/route";

describe("reports JSON API", () => {
  it("rejects unsafe targets without creating a diagnosis", async () => {
    const response = await POST(
      new Request("http://localhost/api/reports", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ target: "localhost" }),
      })
    );
    const body = (await response.json()) as { ok: boolean; outcome?: string };
    expect(response.status).toBe(400);
    expect(body.ok).toBe(false);
    expect(body.outcome).toBe("invalid_input");
  });

  it("creates a shareable unavailable report for a public target", async () => {
    const response = await POST(
      new Request("http://localhost/api/reports", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ target: "youtube.com" }),
      })
    );
    const body = (await response.json()) as {
      ok: boolean;
      document: {
        format: string;
        report: {
          reportId: string;
          likelyCause: string | null;
          insufficientEvidence: { determined: boolean };
          versions: { modelVersion: string };
        };
      };
    };

    expect(response.status).toBe(201);
    expect(body.ok).toBe(true);
    expect(body.document.format).toBe("netpulse.diagnostic-report.v1");
    expect(body.document.report.likelyCause).toBeNull();
    expect(body.document.report.insufficientEvidence.determined).toBe(true);
    expect(body.document.report.versions.modelVersion).toBe("0.5.0");

    const fetched = await GET(new Request("http://localhost/api/reports/x"), {
      params: Promise.resolve({ id: body.document.report.reportId }),
    });
    expect(fetched.status).toBe(200);
  });

  it("does not treat a missing report as a failed diagnosis", async () => {
    const response = await GET(new Request("http://localhost/api/reports/x"), {
      params: Promise.resolve({
        id: "00000000-0000-4000-8000-000000000000",
      }),
    });
    const body = (await response.json()) as { ok: boolean; error: string };
    expect(response.status).toBe(404);
    expect(body.ok).toBe(false);
    expect(body.error).toMatch(/not a measurement failure/i);
  });
});
