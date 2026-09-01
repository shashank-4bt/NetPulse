import { describe, expect, it } from "vitest";

import { startDiagnosis } from "@/features/diagnose/start-diagnosis";
import { getReport } from "@/lib/reports/store";

describe("startDiagnosis", () => {
  it("rejects invalid input without creating a report", async () => {
    const result = await startDiagnosis("localhost");
    expect(result.ok).toBe(false);
    if (!result.ok) {
      expect(result.error).toMatch(/local|internal|private/i);
    }
  });

  it("persists an unavailable report for a valid public target", async () => {
    const result = await startDiagnosis("youtube.com");
    expect(result.ok).toBe(true);
    if (!result.ok) {
      return;
    }
    const report = getReport(result.reportId);
    expect(report?.outcome).toBe("measurement_unavailable");
    expect(report?.likelyCause).toBeNull();
    expect(report?.confidence.percent).toBeNull();
    expect(report?.target.hostname).toBe("youtube.com");
    expect(report?.engineVersion).toBeTruthy();
  });
});
