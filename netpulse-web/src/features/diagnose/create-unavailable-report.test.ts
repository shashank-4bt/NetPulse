import { describe, expect, it } from "vitest";

import {
  CONFIDENCE_CAVEAT,
  EVIDENCE_GRAPH_NODE_IDS,
  INSUFFICIENT_CAUSE_MESSAGE,
} from "@/domain/diagnostic";
import { createUnavailableReport } from "@/features/diagnose/create-unavailable-report";
import { reportHonestyErrors } from "@/features/intelligence/honesty";

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
    expect(report.alternativeHypotheses).toEqual([]);
    expect(report.outcome).toBe("measurement_unavailable");
  });

  it("never marks unmeasured steps or graph nodes as failed", () => {
    const measurementSteps = report.tests.filter(
      (step) => step.id !== "initializing" && step.id !== "complete"
    );
    expect(measurementSteps.every((step) => step.state === "unavailable")).toBe(
      true
    );
    expect(report.tests.some((step) => step.state === "failed")).toBe(false);
    expect(report.graph.map((node) => node.id)).toEqual([
      ...EVIDENCE_GRAPH_NODE_IDS,
    ]);
    expect(report.graph.every((node) => node.status === "not_measured")).toBe(
      true
    );
  });

  it("treats insufficient evidence as a first-class result", () => {
    expect(report.insufficientEvidence.determined).toBe(true);
    expect(report.insufficientEvidence.message).toBe(INSUFFICIENT_CAUSE_MESSAGE);
    expect(report.insufficientEvidence.nextCheck).toMatch(/Next recommended check/i);
  });

  it("versions the engine, rules, measurements, and model", () => {
    expect(report.versions.diagnosticEngineVersion).toBe(report.engineVersion);
    expect(report.versions.ruleVersion).toBeTruthy();
    expect(report.versions.measurementVersion).toBeTruthy();
    expect(report.versions.modelVersion).toBe("0.5.0");
  });

  it("stores structured recommendations that never auto-execute", () => {
    expect(report.recommendations.length).toBeGreaterThan(0);
    for (const recommendation of report.recommendations) {
      expect(recommendation.action).toBeTruthy();
      expect(recommendation.reason).toBeTruthy();
      expect(recommendation.risk).toBeTruthy();
      expect(recommendation.expectedResult).toBeTruthy();
      expect(recommendation.verification).toBeTruthy();
      expect(recommendation.autoExecute).toBe(false);
    }
  });

  it("includes a no-certainty confidence caveat", () => {
    expect(report.confidence.caveat).toBe(CONFIDENCE_CAVEAT);
  });

  it("passes honesty checks", () => {
    expect(reportHonestyErrors(report)).toEqual([]);
  });

  it("does not contain the example hard-coded diagnosis copy", () => {
    const serialized = JSON.stringify(report);
    expect(serialized).not.toMatch(/87%/);
    expect(serialized).not.toMatch(/Likely network\/routing degradation/i);
    expect(serialized).not.toMatch(/TCP failures elevated/i);
  });
});
