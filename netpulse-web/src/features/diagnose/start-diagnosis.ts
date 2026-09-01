"use server";

import { createUnavailableReport } from "@/features/diagnose/create-unavailable-report";
import { validateDiagnoseTarget } from "@/features/diagnose/validate-target";
import { saveReport } from "@/lib/reports/store";

export type StartDiagnosisResult =
  | { ok: true; reportId: string }
  | { ok: false; error: string };

export async function startDiagnosis(raw: string): Promise<StartDiagnosisResult> {
  const validation = validateDiagnoseTarget(raw);
  if (!validation.ok) {
    return { ok: false, error: validation.error };
  }

  const report = createUnavailableReport(validation);
  saveReport(report);
  return { ok: true, reportId: report.reportId };
}
