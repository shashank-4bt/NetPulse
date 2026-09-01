import { afterEach, describe, expect, it, vi } from "vitest";

import { GET } from "@/app/api/map/aggregates/route";

describe("map aggregates BFF", () => {
  afterEach(() => {
    vi.unstubAllEnvs();
    vi.unstubAllGlobals();
  });

  it("returns an empty map when the API is not configured", async () => {
    vi.stubEnv("NETPULSE_API_BASE_URL", "");
    const response = await GET(new Request("http://localhost/api/map/aggregates"));
    const body = (await response.json()) as {
      ok: boolean;
      map: { cells: unknown[] };
    };
    expect(response.status).toBe(503);
    expect(body.ok).toBe(false);
    expect(body.map.cells).toEqual([]);
  });

  it("forwards viewport filters to the API without inventing cells", async () => {
    vi.stubEnv("NETPULSE_API_BASE_URL", "http://api.test");
    const fetchMock = vi.fn(async (input: RequestInfo | URL) => {
      const url = String(input);
      expect(url).toContain("/v1/map/aggregates");
      expect(url).toContain("west=-10");
      return new Response(
        JSON.stringify({
          ok: true,
          map: {
            cells: [],
            incidentRefs: [],
            totalSamples: 0,
            limit: 250,
            truncated: false,
            precision: "none",
            reason: "empty",
            hasCoordinates: false,
          },
        }),
        { status: 200 }
      );
    });
    vi.stubGlobal("fetch", fetchMock);
    const response = await GET(
      new Request(
        "http://localhost/api/map/aggregates?west=-10.4&south=40.2&east=10.4&north=60.2&layers=regional"
      )
    );
    expect(response.status).toBe(200);
    const body = (await response.json()) as { ok: boolean; map: { cells: unknown[] } };
    expect(body.ok).toBe(true);
    expect(body.map.cells).toEqual([]);
  });
});
