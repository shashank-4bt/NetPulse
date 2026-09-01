import { describe, expect, it } from "vitest";

import { createUnavailableReport } from "@/features/diagnose/create-unavailable-report";
import {
  SHAREABLE_AUDIENCES,
  SHAREABLE_FORMAT,
  toShareableReport,
} from "@/features/intelligence/shareable-report";

describe("toShareableReport", () => {
  it("wraps the stored report for web, support, PDF, and JSON API", () => {
    const report = createUnavailableReport({
      ok: true,
      raw: "google.com",
      hostname: "google.com",
      kind: "known_service",
      serviceSlug: "google",
    });
    const document = toShareableReport(report, "2026-09-01T00:00:00.000Z");

    expect(document.format).toBe(SHAREABLE_FORMAT);
    expect(document.audiences).toEqual(SHAREABLE_AUDIENCES);
    expect(document.report.reportId).toBe(report.reportId);
    expect(document.report.versions).toEqual(report.versions);
    expect(document.serializedAt).toBe("2026-09-01T00:00:00.000Z");
  });
});
