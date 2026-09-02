import { describe, expect, it } from "vitest";

import { mapDashboard, mapSLA } from "@/lib/api/developer";

describe("developer honesty mapping", () => {
  it("does not invent percentiles or availability from an empty payload", () => {
    const sla = mapSLA({
      availability: { measured: false, sampleCount: 0, summary: "Observed sample count: 0." },
      downtime: { measured: false, sampleCount: 0 },
      latency: { sampleCount: 0 },
      incidents: [],
      regionalPerformance: [],
      summary: "No SLA window can be computed until monitor checks are stored.",
    });
    expect(sla.availability.measured).toBe(false);
    expect(sla.availability.value).toBeNull();
    expect(sla.latency.p50).toBeNull();
    expect(sla.latency.p95).toBeNull();
    expect(sla.latency.p99).toBeNull();
    expect(sla.summary).not.toMatch(/99\.9|87%/);

    const dash = mapDashboard({
      availability: { measured: false, sampleCount: 0 },
      latency: { sampleCount: 1 },
      incidents: [],
      regionalPerformance: [],
    });
    expect(dash.latency.p95).toBeNull();
    expect(dash.availability.measured).toBe(false);
  });
});
