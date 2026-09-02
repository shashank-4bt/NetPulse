import { describe, expect, it } from "vitest";

import { mapAnalytics, mapDashboard, mapReport } from "@/lib/api/business";

describe("business honesty mapping", () => {
  it("does not invent org health, availability, or percentiles from an empty payload", () => {
    const dash = mapDashboard({
      overallHealth: "Not measured. No stored organization checks exist.",
      availability: { measured: false, sampleCount: 0 },
      incidents: [],
      affectedDevices: [],
      regions: [],
      networks: [],
      services: [],
      summary: "Organization health is computed from stored checks only. Empty series stay unmeasured.",
    });
    expect(dash.availability.measured).toBe(false);
    expect(dash.availability.value).toBeNull();
    expect(dash.overallHealth).not.toMatch(/99\.9|87%/);
    expect(dash.affectedDevices).toHaveLength(0);

    const analytics = mapAnalytics({
      filters: { region: "us-east" },
      availability: { measured: false, sampleCount: 0 },
      latency: { sampleCount: 1 },
      sampleCount: 0,
      incidents: [],
    });
    expect(analytics.latency.p95).toBeNull();
    expect(analytics.availability.measured).toBe(false);
    expect(analytics.filters.region).toBe("us-east");

    const report = mapReport({
      kind: "availability",
      availability: { measured: false, sampleCount: 0 },
      latency: { sampleCount: 0 },
      incidents: [],
      findings: [],
      summary: "No stored checks. This report does not invent availability or latency.",
    });
    expect(report.availability.value).toBeNull();
    expect(report.latency.p50).toBeNull();
    expect(report.summary).not.toMatch(/99\.9|87%/);
  });
});
