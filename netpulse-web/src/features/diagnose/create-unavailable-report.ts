import {
  DIAGNOSTIC_STEP_IDS,
  DIAGNOSTIC_STEP_LABELS,
  DIAGNOSTIC_VERSIONS,
  ENGINE_VERSION,
  EVIDENCE_GRAPH_NODE_IDS,
  EVIDENCE_GRAPH_NODE_LABELS,
  INSUFFICIENT_CAUSE_MESSAGE,
  INSUFFICIENT_NEXT_CHECK,
  MEASUREMENT_BLOCKS,
  MEASUREMENT_LAYER,
  type DiagnosticReport,
  type DiagnosticStep,
  type EvidenceGraphNode,
  type Measurement,
  type Recommendation,
  type VerificationStep,
} from "@/domain/diagnostic";
import { emptyConfidence } from "@/features/intelligence/confidence";
import type { TargetValidation } from "@/features/diagnose/validate-target";

const UNAVAILABLE_NOTE =
  "Measurement workers are not connected. This step was not run and is not a failure.";

function unavailableRecommendation(): Recommendation {
  return {
    id: "re-run-when-workers-exist",
    action: "Re-run this target after measurement workers are connected.",
    reason: "No probes ran, so there is no evidence to act on.",
    risk: "Treating this report as a diagnosis can send users or support down the wrong path.",
    expectedResult: "A later measured run can persist facts before any cause is named.",
    verification:
      "Accept a cause only after a second run records measured facts for the same target.",
    safetyClass: "safe",
    autoExecute: false,
    evidenceIds: [],
  };
}

function unavailableVerification(): VerificationStep {
  return {
    id: "compare-next-measured-run",
    label: "Compare against a later measured run",
    status: "unavailable",
    note: "Verification requires a second measured run. None exists.",
    comparedRunId: null,
  };
}

function unavailableGraph(): EvidenceGraphNode[] {
  return EVIDENCE_GRAPH_NODE_IDS.map((id) => ({
    id,
    label: EVIDENCE_GRAPH_NODE_LABELS[id],
    status: "not_measured",
    confidence: emptyConfidence(),
    timestamp: null,
    evidenceIds: [],
    measurementIds: [],
  }));
}

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
  const measurements: Measurement[] = MEASUREMENT_BLOCKS.map((block) => ({
    id: `measurement-${block.key}`,
    key: block.key,
    label: block.label,
    value: block.key === "timestamp" ? timestamp : null,
    unit: null,
    measured: false,
    measuredAt: null,
    layer: MEASUREMENT_LAYER[block.key] ?? null,
    summary: block.key === "timestamp" ? timestamp : null,
  }));

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
    measurements,
    evidence: [],
    hypotheses: [],
    alternativeHypotheses: [],
    likelyCause: null,
    confidence: emptyConfidence(),
    recommendations: [unavailableRecommendation()],
    verificationSteps: [unavailableVerification()],
    escalationConditions: [
      {
        id: "do-not-escalate-without-facts",
        when: "This report has no measured evidence.",
        action:
          "Do not escalate to an ISP, CDN, or vendor using this report as proof of a failure.",
        safetyClass: "advisory",
      },
    ],
    graph: unavailableGraph(),
    versions: DIAGNOSTIC_VERSIONS,
    insufficientEvidence: {
      determined: true,
      message: INSUFFICIENT_CAUSE_MESSAGE,
      nextCheck: INSUFFICIENT_NEXT_CHECK,
    },
    engineVersion: ENGINE_VERSION,
  };
}
