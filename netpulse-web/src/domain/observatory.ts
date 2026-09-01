import type {
  Confidence,
  Evidence,
  Hypothesis,
} from "@/domain/diagnostic";
import type { LayerStatus, Severity } from "@/domain/display";

export const INCIDENT_STAGES = [
  "detected",
  "investigating",
  "correlated",
  "identified",
  "mitigating",
  "resolved",
] as const;

export type IncidentStage = (typeof INCIDENT_STAGES)[number];

export const INCIDENT_STAGE_LABELS: Record<IncidentStage, string> = {
  detected: "Detected",
  investigating: "Investigating",
  correlated: "Correlated",
  identified: "Identified",
  mitigating: "Mitigating",
  resolved: "Resolved",
};

export const TIMELINE_EVENT_STATES = [
  "complete",
  "current",
  "not_reached",
] as const;

export type TimelineEventState = (typeof TIMELINE_EVENT_STATES)[number];

export const SERVICE_CURRENT_STATES = [
  "not_measured",
  "insufficient_evidence",
  "unknown",
  "operational",
  "degraded",
  "investigating",
  "major_incident",
] as const;

export type ServiceCurrentState = (typeof SERVICE_CURRENT_STATES)[number];

export type MetricObservation = {
  value: string | number | null;
  unit: string | null;
  measured: boolean;
  sampleCount: number;
  sampleWindow: string | null;
  summary: string;
};

export type SliceObservation = {
  id: string;
  label: string;
  status: LayerStatus;
  sampleCount: number;
  summary: string;
};

export type ServiceIntelligence = {
  currentState: ServiceCurrentState;
  health: number | null;
  lastUpdated: string | null;
  availability: MetricObservation;
  latency: MetricObservation;
  errors: MetricObservation;
  regionalHealth: SliceObservation[];
  networkHealth: SliceObservation[];
  recentIncidentIds: string[];
};

export type IncidentTimelineEvent = {
  stage: IncidentStage;
  label: string;
  status: TimelineEventState;
  at: string | null;
  note: string;
};

export type PublicIncidentRecord = {
  id: string;
  title: string;
  severity: Severity | "";
  status: IncidentStage | "";
  scope: string;
  startedAt: string;
  lastUpdatedAt: string;
  affectedServices: string[];
  regions: string[];
  networks: string[];
  evidence: Evidence[];
  hypotheses: Hypothesis[];
  confidence: Confidence;
  timeline: IncidentTimelineEvent[];
  sampleCount: number;
  sampleRate: string | null;
  affectedUserCount: null;
};

export type OutageQuery = {
  service: string;
  region: string;
  network: string;
  severity: string;
  status: string;
  time: string;
  q: string;
  sort: string;
  page: number;
};

export const OUTAGE_TIME_WINDOWS = ["all", "24h", "7d", "30d"] as const;
export const OUTAGE_SORTS = [
  "started_desc",
  "started_asc",
  "updated_desc",
  "severity",
  "status",
] as const;

export const DEFAULT_OUTAGE_QUERY: OutageQuery = {
  service: "",
  region: "",
  network: "",
  severity: "",
  status: "",
  time: "all",
  q: "",
  sort: "started_desc",
  page: 1,
};

export const PAGE_SIZE = 20;
export const OBSERVED_FAILURES_TITLE = "Elevated connectivity failures observed";
