import { describe, expect, it } from "vitest";

import { CONFIDENCE_CAVEAT } from "@/domain/diagnostic";
import type { PublicIncidentRecord } from "@/domain/observatory";
import { timelineFor } from "@/features/observatory/timeline";

describe("incident timeline", () => {
  it("does not mark resolved from a single recovery sample", () => {
    const incident: PublicIncidentRecord = {
      id: "11111111-1111-4111-8111-111111111111",
      title: "Elevated connectivity failures observed",
      severity: "high",
      status: "resolved",
      scope: "youtube",
      startedAt: "2026-09-01T04:00:00Z",
      lastUpdatedAt: "2026-09-01T05:00:00Z",
      affectedServices: ["youtube"],
      regions: [],
      networks: [],
      evidence: [
        {
          id: "e1",
          evidenceClass: "measured_fact",
          title: "Worker HTTP observation",
          body: "status=200",
          measurementIds: ["m1"],
          layer: "service",
          observedAt: "2026-09-01T05:00:00Z",
        },
      ],
      hypotheses: [],
      confidence: {
        level: null,
        percent: null,
        supportingEvidenceIds: [],
        alternativeHypothesisIds: [],
        caveat: CONFIDENCE_CAVEAT,
      },
      timeline: [],
      sampleCount: 1,
      sampleRate: "1 worker probe",
      affectedUserCount: null,
    };
    const resolved = timelineFor(incident).find((event) => event.stage === "resolved");
    expect(resolved?.status).toBe("not_reached");
  });
});
