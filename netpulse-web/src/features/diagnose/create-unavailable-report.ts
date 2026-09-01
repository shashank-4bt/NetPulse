import {
  DIAGNOSTIC_STEP_IDS,
  DIAGNOSTIC_STEP_LABELS,
  ENGINE_VERSION,
  MEASUREMENT_BLOCKS,
  type DiagnosticReport,
  type DiagnosticStep,
} from "@/domain/diagnostic";
import type { TargetValidation } from "@/features/diagnose/validate-target";

const UNAVAILABLE_NOTE =
  "Measurement workers are not connected. This step was not run and is not a failure.";

export function createUnavailableReport(
  validation: Extract<TargetValidation, { ok: true }>
): DiagnosticReport {
  const tests: DiagnosticStep[] = DIAGNOSTIC_STEP_IDS.map((id) => {
    if (id === "initializing" || id === "complete") {
      return {
        id,
        label: DIAGNOSTIC_STEP_LABELS[id],
        state: "complete",
        durationMs: null,
        note:
          id === "initializing"
            ? "Run accepted. No worker was scheduled."
            : "Run finished without measurements.",
      };
    }
    return {
      id,
      label: DIAGNOSTIC_STEP_LABELS[id],
      state: "unavailable",
      durationMs: null,
      note: UNAVAILABLE_NOTE,
    };
  });

  const timestamp = new Date().toISOString();

  return {
    reportId: crypto.randomUUID(),
    target: {
      raw: validation.raw,
      hostname: validation.hostname,
      kind: validation.kind,
      serviceSlug: validation.serviceSlug,
    },
    timestamp,
    outcome: "measurement_unavailable",
    tests,
    measurements: MEASUREMENT_BLOCKS.map((block) => ({
      key: block.key,
      label: block.label,
      summary: block.key === "timestamp" ? timestamp : null,
      measured: false,
    })),
    evidence: [],
    hypotheses: [],
    likelyCause: null,
    confidence: {
      level: null,
      percent: null,
    },
    recommendation:
      "Do not treat this page as a diagnosis. Re-run when measurement workers are connected.",
    verification: {
      status: "unavailable",
      note: "Verification requires a second measured run. None exists.",
    },
    engineVersion: ENGINE_VERSION,
  };
}
