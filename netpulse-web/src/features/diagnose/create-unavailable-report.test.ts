import { describe, expect, it } from "vitest";

import { createUnavailableReport } from "@/features/diagnose/create-unavailable-report";

describe("createUnavailableReport", () => {
  const report = createUnavailableReport({
    ok: true,
    raw: "youtube.com",
    hostname: "youtube.com",
    kind: "known_service",
    serviceSlug: "youtube",
  });

  it("does not invent a likely cause, confidence percent, or evidence", () => {
    expect(report.likelyCause).toBeNull();
    expect(report.confidence.percent).toBeNull();
    expect(report.confidence.level).toBeNull();
    expect(report.evidence).toEqual([]);
    expect(report.hypotheses).toEqual([]);
    expect(report.outcome).toBe("measurement_unavailable");
  });

  it("never marks unmeasured steps as failed", () => {
    const measurementSteps = report.tests.filter(
      (step) => step.id !== "initializing" && step.id !== "complete"
    );
    expect(measurementSteps.every((step) => step.state === "unavailable")).toBe(
      true
    );
    expect(report.tests.some((step) => step.state === "failed")).toBe(false);
  });

  it("does not contain the example hard-coded diagnosis copy", () => {
    const serialized = JSON.stringify(report);
    expect(serialized).not.toMatch(/87%/);
    expect(serialized).not.toMatch(/Likely network\/routing degradation/i);
    expect(serialized).not.toMatch(/TCP failures elevated/i);
  });
});
