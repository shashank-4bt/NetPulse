import { describe, expect, it } from "vitest";

import { mapSystem } from "@/lib/api/admin";

describe("admin honesty mapping", () => {
  it("does not invent error rates or latency percentiles from an empty payload", () => {
    const system = mapSystem({
      api: { name: "API", status: "up", detail: "This process answered the request.", measured: true },
      worker: { name: "Worker", status: "unmeasured", detail: "No worker heartbeat is stored.", measured: false },
      queue: { name: "Queue", status: "observed", detail: "Queue depth: 0.", measured: true },
      database: { name: "Database", status: "configured", detail: "Database adapter: memory.", measured: true },
      cache: { name: "Cache", status: "configured", detail: "Cache adapter: memory.", measured: true },
      measurementFailures: { measured: false, sampleCount: 0 },
      errorRates: { measured: false, sampleCount: 0 },
      latency: { sampleCount: 0 },
      summary: "System figures use stored process state and stored measurements only. Empty series stay unmeasured.",
    });
    expect(system.errorRates.measured).toBe(false);
    expect(system.errorRates.value).toBeNull();
    expect(system.measurementFailures.measured).toBe(false);
    expect(system.latency.p50).toBeNull();
    expect(system.latency.p95).toBeNull();
    expect(system.latency.p99).toBeNull();
    expect(system.worker.measured).toBe(false);
    expect(system.summary).not.toMatch(/99\.9|87%/);
  });
});
