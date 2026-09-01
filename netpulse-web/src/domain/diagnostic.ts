import type { ConfidenceLevel, EvidenceClass } from "@/domain/display";

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

export type DiagnosticEvidence = {
  id: string;
  evidenceClass: EvidenceClass;
  title: string;
  body: string;
};

export type DiagnosticHypothesis = {
  id: string;
  title: string;
  body: string;
};

export type DiagnosticConfidence = {
  level: ConfidenceLevel | null;
  percent: number | null;
};

export type DiagnosticMeasurementBlock = {
  key: string;
  label: string;
  summary: string | null;
  measured: boolean;
};

export type DiagnosticReport = {
  reportId: string;
  target: DiagnosticTarget;
  timestamp: string;
  outcome: DiagnosticOutcome;
  tests: DiagnosticStep[];
  measurements: DiagnosticMeasurementBlock[];
  evidence: DiagnosticEvidence[];
  hypotheses: DiagnosticHypothesis[];
  likelyCause: string | null;
  confidence: DiagnosticConfidence;
  recommendation: string | null;
  verification: {
    status: "not_run" | "unavailable";
    note: string;
  };
  engineVersion: string;
};

export const ENGINE_VERSION = "0.4.0-unavailable";

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

export const MEASUREMENT_BLOCKS: readonly { key: string; label: string }[] = [
  { key: "dns", label: "DNS" },
  { key: "tcp", label: "TCP" },
  { key: "tls", label: "TLS" },
  { key: "http", label: "HTTP" },
  { key: "latency", label: "Latency" },
  { key: "packet_loss", label: "Packet loss" },
  { key: "routing", label: "Routing" },
  { key: "network", label: "Network" },
  { key: "region", label: "Region" },
  { key: "timestamp", label: "Timestamp" },
] as const;
