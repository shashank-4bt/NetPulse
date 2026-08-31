/**
 * Presentation taxonomies for NetPulse UI.
 * These are display contracts only — they do not compute scores or invent evidence.
 */

export const OPERATIONAL_STATUSES = [
  "operational",
  "degraded",
  "investigating",
  "major_incident",
  "unknown",
] as const;

export type OperationalStatus = (typeof OPERATIONAL_STATUSES)[number];

export const SEVERITIES = [
  "informational",
  "moderate",
  "high",
  "critical",
] as const;

export type Severity = (typeof SEVERITIES)[number];

export const CONFIDENCE_LEVELS = [
  "low",
  "medium",
  "high",
  "very_high",
] as const;

export type ConfidenceLevel = (typeof CONFIDENCE_LEVELS)[number];

export const EVIDENCE_CLASSES = [
  "measured_fact",
  "inferred_hypothesis",
  "recommendation",
] as const;

export type EvidenceClass = (typeof EVIDENCE_CLASSES)[number];

export const LAYER_STATUSES = [
  "not_measured",
  "insufficient_evidence",
  "healthy",
  "degraded",
  "failed",
] as const;

export type LayerStatus = (typeof LAYER_STATUSES)[number];
