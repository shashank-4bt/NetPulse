import { describe, expect, it } from "vitest";

import { createUnavailableReport } from "@/features/diagnose/create-unavailable-report";
import { parseReportId } from "@/lib/reports/id";
import { getReport, saveReport } from "@/lib/reports/store";

describe("report store", () => {
  it("saves and returns a report by id", () => {
    const report = createUnavailableReport({
      ok: true,
      raw: "example.com",
      hostname: "example.com",
      kind: "domain",
      serviceSlug: null,
    });
    saveReport(report);
    expect(getReport(report.reportId)?.reportId).toBe(report.reportId);
  });

  it("rejects malformed ids instead of treating them as a failed diagnosis", () => {
    expect(parseReportId("../etc/passwd")).toBeNull();
    expect(parseReportId("not-a-uuid")).toBeNull();
    expect(getReport("not-a-uuid")).toBeNull();
  });
});
