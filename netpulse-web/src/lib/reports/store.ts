import type { DiagnosticReport } from "@/domain/diagnostic";
import { parseReportId } from "@/lib/reports/id";

/**
 * Process-local report store.
 * PRODUCTION ENGINEERING REQUIREMENT until PostgreSQL exists.
 * Reports do not survive process restart and are not shared across instances.
 *
 * The Map is hung on globalThis so App Router, server actions, and route
 * handlers share one store instead of one Map per bundled module instance.
 */
const globalForReports = globalThis as typeof globalThis & {
  __netpulseReports?: Map<string, DiagnosticReport>;
};

const reports =
  globalForReports.__netpulseReports ?? new Map<string, DiagnosticReport>();

globalForReports.__netpulseReports = reports;

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
