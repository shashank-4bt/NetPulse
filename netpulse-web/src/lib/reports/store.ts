import type { DiagnosticReport } from "@/domain/diagnostic";
import { parseReportId } from "@/lib/reports/id";

/**
 * Process-local report store.
 * PRODUCTION ENGINEERING REQUIREMENT until PostgreSQL exists.
 * Reports do not survive process restart and are not shared across instances.
 */
const reports = new Map<string, DiagnosticReport>();

export function saveReport(report: DiagnosticReport): DiagnosticReport {
  const id = parseReportId(report.reportId);
  if (!id) {
    throw new Error("Report id is not a UUID.");
  }
  const normalized = { ...report, reportId: id };
  reports.set(id, normalized);
  return normalized;
}

export function getReport(reportId: string): DiagnosticReport | null {
  const id = parseReportId(reportId);
  if (!id) {
    return null;
  }
  return reports.get(id) ?? null;
}
