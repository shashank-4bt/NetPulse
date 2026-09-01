"use server";

import { createUnavailableReport } from "@/features/diagnose/create-unavailable-report";
import { validateDiagnoseTarget } from "@/features/diagnose/validate-target";
import {
  getDiagnosis,
  isApiConfigured,
  postDiagnosis,
} from "@/lib/api/backend";
import { saveReport } from "@/lib/reports/store";

export type StartDiagnosisResult =
  | { ok: true; reportId: string }
  | { ok: false; error: string };

export async function startDiagnosis(raw: string): Promise<StartDiagnosisResult> {
  const validation = validateDiagnoseTarget(raw);
  if (!validation.ok) {
    return { ok: false, error: validation.error };
  }

  if (!isApiConfigured()) {
    const report = createUnavailableReport(validation);
    saveReport(report);
    return { ok: true, reportId: report.reportId };
  }

  const created = await postDiagnosis(raw);
  if (!created.ok) {
    return { ok: false, error: created.message };
  }

  const deadline = Date.now() + 12000;
  while (Date.now() < deadline) {
    const latest = await getDiagnosis(created.diagnosis.id);
    if (latest.ok && latest.diagnosis.report) {
      saveReport(latest.diagnosis.report);
      return { ok: true, reportId: created.diagnosis.id };
    }
    if (latest.ok && isTerminal(latest.diagnosis.status) && !latest.diagnosis.report) {
      return { ok: true, reportId: created.diagnosis.id };
    }
    await wait(300);
  }
  return { ok: true, reportId: created.diagnosis.id };
}

function isTerminal(status: string): boolean {
  return (
    status === "complete" ||
    status === "partial" ||
    status === "failed" ||
    status === "unavailable" ||
    status === "insufficient_evidence"
  );
}

function wait(ms: number): Promise<void> {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}
