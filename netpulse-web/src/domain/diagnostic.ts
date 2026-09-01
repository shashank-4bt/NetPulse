import type { ConfidenceLevel, EvidenceClass, LayerStatus } from "@/domain/display";

export const DIAGNOSTIC_STEP_IDS = [
  "initializing",
  "device",
  "wifi",
  "dns",
  "connectivity",
  "isp",
  "routing",
  "tls",
  "http",
  "service",
  "regional_comparison",
  "network_comparison",
  "analysis",
  "complete",
] as const;

export type DiagnosticStepId = (typeof DIAGNOSTIC_STEP_IDS)[number];

export const DIAGNOSTIC_STEP_STATES = [
  "pending",
  "current",
  "complete",
  "failed",
  "unavailable",
  "not_measured",
] as const;

export type DiagnosticStepState = (typeof DIAGNOSTIC_STEP_STATES)[number];

export const DIAGNOSTIC_OUTCOMES = [
  "success",
  "partial_success",
  "timeout",
  "backend_unavailable",
  "insufficient_evidence",
  "invalid_input",
  "measurement_unavailable",
] as const;

export type DiagnosticOutcome = (typeof DIAGNOSTIC_OUTCOMES)[number];

export type TargetKind = "domain" | "url" | "known_service";

export type DiagnosticTarget = {
  raw: string;
  hostname: string;
  kind: TargetKind;
  serviceSlug: string | null;
};

export type DiagnosticStep = {
  id: DiagnosticStepId;
  label: string;
  state: DiagnosticStepState;
  durationMs: number | null;
  note: string | null;
};

export const MEASUREMENT_KEYS = [
  "dns",
  "tcp",
  "tls",
  "http",
  "latency",
  "packet_loss",
  "routing",
  "network",
  "region",
  "timestamp",
] as const;

export type MeasurementKey = (typeof MEASUREMENT_KEYS)[number];

export const EVIDENCE_GRAPH_NODE_IDS = [
  "device",
  "wifi",
  "router",
  "isp",
  "route",
  "cdn",
  "service",
] as const;

export type EvidenceGraphNodeId = (typeof EVIDENCE_GRAPH_NODE_IDS)[number];

export const EVIDENCE_GRAPH_NODE_LABELS: Record<EvidenceGraphNodeId, string> = {
  device: "Device",
  wifi: "Wi-Fi",
  router: "Router",
  isp: "ISP",
  route: "Route",
  cdn: "CDN",
  service: "Service",
};

export const RECOMMENDATION_SAFETY_CLASSES = [
  "safe",
  "advisory",
  "dangerous",
] as const;

export type RecommendationSafetyClass =
  (typeof RECOMMENDATION_SAFETY_CLASSES)[number];

export type Measurement = {
  id: string;
  key: MeasurementKey;
  label: string;
  value: string | number | null;
  unit: string | null;
  measured: boolean;
  measuredAt: string | null;
  layer: EvidenceGraphNodeId | null;
  summary: string | null;
};

export type Evidence = {
  id: string;
  evidenceClass: EvidenceClass;
  title: string;
  body: string;
  measurementIds: string[];
  layer: EvidenceGraphNodeId | null;
  observedAt: string | null;
};

export type Confidence = {
  level: ConfidenceLevel | null;
  percent: number | null;
  supportingEvidenceIds: string[];
  alternativeHypothesisIds: string[];
  caveat: string;
};

export type Hypothesis = {
  id: string;
  title: string;
  body: string;
  evidenceIds: string[];
  layer: EvidenceGraphNodeId | null;
  confidence: Confidence;
};

export type AlternativeHypothesis = Hypothesis & {
  kind: "alternative";
};

export type Recommendation = {
  id: string;
  action: string;
  reason: string;
  risk: string;
  expectedResult: string;
  verification: string;
  safetyClass: RecommendationSafetyClass;
  autoExecute: false;
  evidenceIds: string[];
};

export type VerificationStep = {
  id: string;
  label: string;
  status: "not_run" | "unavailable" | "passed" | "failed";
  note: string;
  comparedRunId: string | null;
};

export type EscalationCondition = {
  id: string;
  when: string;
  action: string;
  safetyClass: Extract<RecommendationSafetyClass, "advisory" | "dangerous">;
};

export type DiagnosticVersions = {
  diagnosticEngineVersion: string;
  ruleVersion: string;
  measurementVersion: string;
  modelVersion: string;
};

export type EvidenceGraphNode = {
  id: EvidenceGraphNodeId;
  label: string;
  status: LayerStatus;
  confidence: Confidence;
  timestamp: string | null;
  evidenceIds: string[];
  measurementIds: string[];
};

export type InsufficientEvidenceResult = {
  determined: boolean;
  message: string;
  nextCheck: string;
};

export type DiagnosticReport = {
  reportId: string;
  target: DiagnosticTarget;
  timestamp: string;
  outcome: DiagnosticOutcome;
  tests: DiagnosticStep[];
  measurements: Measurement[];
  evidence: Evidence[];
  hypotheses: Hypothesis[];
  alternativeHypotheses: AlternativeHypothesis[];
  likelyCause: string | null;
  confidence: Confidence;
  recommendations: Recommendation[];
  verificationSteps: VerificationStep[];
  escalationConditions: EscalationCondition[];
  graph: EvidenceGraphNode[];
  versions: DiagnosticVersions;
  insufficientEvidence: InsufficientEvidenceResult;
  engineVersion: string;
};

export const MODEL_VERSION = "0.5.0";
export const ENGINE_VERSION = "0.5.0-unavailable";
export const RULE_VERSION = "0.0.0-none";
export const MEASUREMENT_VERSION = "0.0.0-none";

export const DIAGNOSTIC_VERSIONS: DiagnosticVersions = {
  diagnosticEngineVersion: ENGINE_VERSION,
  ruleVersion: RULE_VERSION,
  measurementVersion: MEASUREMENT_VERSION,
  modelVersion: MODEL_VERSION,
};

export const CONFIDENCE_CAVEAT =
  "Confidence is not certainty. A level or percentage can be wrong and must not be treated as proof.";

export const INSUFFICIENT_CAUSE_MESSAGE =
  "NetPulse could not safely determine the root cause.";

export const INSUFFICIENT_NEXT_CHECK =
  "Next recommended check: connect measurement workers and re-run this target. Do not change device, DNS, or ISP settings based on this report.";

export const DIAGNOSTIC_STEP_LABELS: Record<DiagnosticStepId, string> = {
  initializing: "Initializing",
  device: "Device",
  wifi: "Wi-Fi",
  dns: "DNS",
  connectivity: "Connectivity",
  isp: "ISP",
  routing: "Routing",
  tls: "TLS",
  http: "HTTP",
  service: "Service",
  regional_comparison: "Regional comparison",
  network_comparison: "Network comparison",
  analysis: "Analysis",
  complete: "Complete",
};

export const MEASUREMENT_BLOCKS: readonly {
  key: MeasurementKey;
  label: string;
}[] = [
  { key: "dns", label: "DNS" },
  { key: "tcp", label: "TCP" },
  { key: "tls", label: "TLS" },
  { key: "http", label: "HTTP" },
  { key: "latency", label: "Latency" },
  { key: "packet_loss", label: "Packet loss" },
  { key: "routing", label: "Route metadata" },
  { key: "network", label: "Network/ASN" },
  { key: "region", label: "Region" },
  { key: "timestamp", label: "Timestamp" },
] as const;

export const MEASUREMENT_LAYER: Partial<
  Record<MeasurementKey, EvidenceGraphNodeId>
> = {
  routing: "route",
  network: "isp",
};
