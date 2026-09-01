import { afterEach, describe, expect, it, vi } from "vitest";

import { emptyMapAggregates } from "@/domain/map";
import { loadMapAggregates } from "@/lib/map/load";

describe("loadMapAggregates", () => {
  afterEach(() => {
    vi.unstubAllEnvs();
    vi.unstubAllGlobals();
  });

  it("stays unavailable without an API and does not invent cells", async () => {
    vi.stubEnv("NETPULSE_API_BASE_URL", "");
    const loaded = await loadMapAggregates({});
    expect(loaded.state).toBe("unavailable");
    expect(loaded.aggregates.cells).toEqual([]);
    expect(loaded.aggregates.hasCoordinates).toBe(false);
  });

  it("maps an empty API payload as empty, not healthy", async () => {
    vi.stubEnv("NETPULSE_API_BASE_URL", "http://api.test");
    vi.stubGlobal(
      "fetch",
      vi.fn(async () =>
        new Response(
          JSON.stringify({
            ok: true,
            map: emptyMapAggregates("No coarse geographic aggregates are stored."),
          }),
          { status: 200 }
        )
      )
    );
    const loaded = await loadMapAggregates({ layers: "regional" });
    expect(loaded.state).toBe("empty");
    expect(loaded.aggregates.cells).toEqual([]);
  });

  it("rejects payloads that include raw measurements", async () => {
    vi.stubEnv("NETPULSE_API_BASE_URL", "http://api.test");
    vi.stubGlobal(
      "fetch",
      vi.fn(async () =>
        new Response(
          JSON.stringify({
            ok: true,
            map: {
              cells: [],
              incidentRefs: [],
              measurements: [{ id: "raw" }],
              totalSamples: 0,
              limit: 250,
              truncated: false,
              precision: "none",
              reason: "bad",
              hasCoordinates: false,
            },
          }),
          { status: 200 }
        )
      )
    );
    const loaded = await loadMapAggregates({});
    expect(loaded.state).toBe("error");
    expect(loaded.aggregates.cells).toEqual([]);
  });

  it("strips precise coordinates from cells", async () => {
    vi.stubEnv("NETPULSE_API_BASE_URL", "http://api.test");
    vi.stubGlobal(
      "fetch",
      vi.fn(async () =>
        new Response(
          JSON.stringify({
            ok: true,
            map: {
              level: "world",
              cells: [
                {
                  id: "fixture:london",
                  level: "region",
                  label: "london",
                  lon: -0.12,
                  lat: 51.5,
                  sampleCount: 2,
                  status: "insufficient_evidence",
                  summary: "fixture",
                  layer: "regional",
                  childCount: 0,
                },
              ],
              incidentRefs: [],
              totalSamples: 2,
              limit: 250,
              truncated: false,
              precision: "degree",
              reason: "fixture",
              hasCoordinates: true,
            },
          }),
          { status: 200 }
        )
      )
    );
    const loaded = await loadMapAggregates({});
    expect(loaded.aggregates.cells[0]?.lon).toBeNull();
    expect(loaded.aggregates.cells[0]?.lat).toBeNull();
  });
});
