import { describe, expect, it } from "vitest";

import { CONFIDENCE_CAVEAT } from "@/domain/diagnostic";
import type { PublicIncidentRecord } from "@/domain/observatory";
import { filterIncidents } from "@/features/observatory/filter-incidents";

function fixture(partial: Partial<PublicIncidentRecord>): PublicIncidentRecord {
  return {
    id: "11111111-1111-4111-8111-111111111111",
    title: "Elevated connectivity failures observed",
    severity: "high",
    status: "investigating",
    scope: "youtube",
    startedAt: "2026-09-01T04:00:00Z",
    lastUpdatedAt: "2026-09-01T05:00:00Z",
    affectedServices: ["youtube"],
    regions: ["eu-west"],
    networks: ["AS64500"],
    evidence: [],
    hypotheses: [],
    confidence: {
      level: null,
      percent: null,
      supportingEvidenceIds: [],
      alternativeHypothesisIds: [],
      caveat: CONFIDENCE_CAVEAT,
    },
    timeline: [],
    sampleCount: 3,
    sampleRate: "3 worker probes / 15 min",
    affectedUserCount: null,
    ...partial,
  };
}

describe("filterIncidents", () => {
  it("does not invent rows from an empty store", () => {
    const result = filterIncidents([], {
      service: "",
      region: "",
      network: "",
      severity: "",
      status: "",
      time: "all",
      q: "",
      sort: "started_desc",
      page: 1,
    });
    expect(result.total).toBe(0);
    expect(result.items).toEqual([]);
  });

  it("filters, searches, and paginates stored records only", () => {
    const now = Date.parse("2026-09-01T12:00:00Z");
    const items = [
      fixture({}),
      fixture({
        id: "22222222-2222-4222-8222-222222222222",
        scope: "github",
        affectedServices: ["github"],
        regions: ["us-east"],
        networks: ["AS64501"],
        startedAt: "2026-08-01T00:00:00Z",
        severity: "moderate",
        status: "detected",
      }),
    ];
    const result = filterIncidents(
      items,
      {
        service: "youtube",
        region: "eu-west",
        network: "AS64500",
        severity: "high",
        status: "investigating",
        time: "30d",
        q: "elevated",
        sort: "started_desc",
        page: 1,
      },
      now
    );
    expect(result.total).toBe(1);
    expect(result.items[0]?.id).toBe("11111111-1111-4111-8111-111111111111");
    expect(result.items[0]?.affectedUserCount).toBeNull();
  });
});
