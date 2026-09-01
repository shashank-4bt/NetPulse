import {
  CONFIDENCE_CAVEAT,
  type Confidence,
} from "@/domain/diagnostic";

export function emptyConfidence(): Confidence {
  return {
    level: null,
    percent: null,
    supportingEvidenceIds: [],
    alternativeHypothesisIds: [],
    caveat: CONFIDENCE_CAVEAT,
  };
}

export function formatConfidenceValue(confidence: Confidence): string {
  if (confidence.percent === null) {
    return "No numeric value";
  }
  return `${confidence.percent}%`;
}
