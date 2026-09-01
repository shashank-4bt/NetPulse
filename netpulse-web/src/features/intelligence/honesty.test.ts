import { describe, expect, it } from "vitest";

import { createUnavailableReport } from "@/features/diagnose/create-unavailable-report";
import { reportHonestyErrors } from "@/features/intelligence/honesty";

describe("reportHonestyErrors", () => {
  it("rejects converting an unmeasured node into failure", () => {
    const report = createUnavailableReport({
      ok: true,
      raw: "example.com",
      hostname: "example.com",
      kind: "domain",
      serviceSlug: null,
    });
    const first = report.graph[0];
    if (!first) {
      throw new Error("expected a graph node");
    }
    report.graph[0] = {
      ...first,
      status: "failed",
    };
    expect(reportHonestyErrors(report)).toContain(
      "Unknown or unmeasured nodes must not be marked failed."
    );
  });

  it("rejects a likely cause when no facts exist", () => {
    const report = createUnavailableReport({
      ok: true,
      raw: "example.com",
      hostname: "example.com",
      kind: "domain",
      serviceSlug: null,
    });
    report.likelyCause = "Likely network/routing degradation";
    expect(reportHonestyErrors(report)).toContain(
      "A likely cause requires measured facts."
    );
  });
});
