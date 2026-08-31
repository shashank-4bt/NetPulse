import type {
  ConfidenceLevel,
  EvidenceClass,
  OperationalStatus,
  Severity,
} from "@/domain/display";

export const STATUS_LABEL: Record<OperationalStatus, string> = {
  operational: "Operational",
  degraded: "Degraded",
  investigating: "Investigating",
  major_incident: "Major Incident",
  unknown: "Unknown",
};

export const SEVERITY_LABEL: Record<Severity, string> = {
  informational: "Informational",
  moderate: "Moderate",
  high: "High",
  critical: "Critical",
};

export const CONFIDENCE_LABEL: Record<ConfidenceLevel, string> = {
  low: "Low",
  medium: "Medium",
  high: "High",
  very_high: "Very High",
};

export const EVIDENCE_CLASS_LABEL: Record<EvidenceClass, string> = {
  measured_fact: "Measured fact",
  inferred_hypothesis: "Inferred hypothesis",
  recommendation: "Recommendation",
};

export const INSUFFICIENT_EVIDENCE_LABEL = "Insufficient evidence";
