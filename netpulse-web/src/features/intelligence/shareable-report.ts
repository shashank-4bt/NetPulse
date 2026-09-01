import type { DiagnosticReport } from "@/domain/diagnostic";

export const SHAREABLE_FORMAT = "netpulse.diagnostic-report.v1";

export const SHAREABLE_AUDIENCES = [
  "web",
  "support",
  "pdf",
  "json-api",
] as const;

export type ShareableAudience = (typeof SHAREABLE_AUDIENCES)[number];

export type ShareableDiagnosticReport = {
  format: typeof SHAREABLE_FORMAT;
  formatVersion: "1";
  audiences: typeof SHAREABLE_AUDIENCES;
  serializedAt: string;
  report: DiagnosticReport;
};

export function toShareableReport(
  report: DiagnosticReport,
  serializedAt = new Date().toISOString()
): ShareableDiagnosticReport {
  return {
    format: SHAREABLE_FORMAT,
    formatVersion: "1",
    audiences: SHAREABLE_AUDIENCES,
    serializedAt,
    report,
  };
}

export function stringifyShareableReport(report: DiagnosticReport): string {
  return `${JSON.stringify(toShareableReport(report), null, 2)}\n`;
}
