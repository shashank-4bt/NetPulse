import { getDiagnosis, isApiConfigured } from "@/lib/api/backend";
import type { DiagnosticReport } from "@/domain/diagnostic";
import { getReport, saveReport } from "@/lib/reports/store";

export type LoadedDiagnosis = {
  report: DiagnosticReport | null;
  status: string | null;
  missing: boolean;
  pending: boolean;
  backendError: string | null;
};

export async function loadDiagnosis(id: string): Promise<LoadedDiagnosis> {
  if (isApiConfigured()) {
    const result = await getDiagnosis(id);
    if (!result.ok) {
      if (result.code === "not_found") {
        return {
          report: null,
          status: null,
          missing: true,
          pending: false,
          backendError: null,
        };
      }
      return {
        report: null,
        status: null,
        missing: false,
        pending: false,
        backendError: result.message,
      };
    }
    if (result.diagnosis.report) {
      saveReport(result.diagnosis.report);
    }
    const pending =
      result.diagnosis.status === "queued" ||
      result.diagnosis.status === "running";
    return {
      report: result.diagnosis.report,
      status: result.diagnosis.status,
      missing: false,
      pending,
      backendError: null,
    };
  }

  const report = getReport(id);
  return {
    report,
    status: report ? report.outcome : null,
    missing: !report,
    pending: false,
    backendError: null,
  };
}
