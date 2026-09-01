import { CONFIDENCE_CAVEAT } from "@/domain/diagnostic";
import type { ServiceCatalogEntry } from "@/lib/content/services";
import type {
  MetricObservation,
  PublicIncidentRecord,
  ServiceIntelligence,
} from "@/domain/observatory";

export function unmeasuredObservation(topic: string): MetricObservation {
  return {
    value: null,
    unit: null,
    measured: false,
    sampleCount: 0,
    sampleWindow: null,
    summary: `Observed sample count: 0. ${topic} is not estimated for a population.`,
  };
}

export function emptyServiceIntelligence(): ServiceIntelligence {
  return {
    currentState: "not_measured",
    health: null,
    lastUpdated: null,
    availability: unmeasuredObservation("Availability"),
    latency: unmeasuredObservation("Latency"),
    errors: unmeasuredObservation("Error rate"),
    regionalHealth: [],
    networkHealth: [],
    recentIncidentIds: [],
  };
}

export function sampleCaption(sampleCount: number): string {
  return `Observed sample count: ${sampleCount}. Population impact is not estimated.`;
}

export function catalogWithoutLiveStatus(entry: ServiceCatalogEntry): {
  slug: string;
  name: string;
  category: string;
  summary: string;
  layers: readonly string[];
} {
  return {
    slug: entry.slug,
    name: entry.name,
    category: entry.category,
    summary: entry.summary,
    layers: entry.layers,
  };
}

export function emptyIncidentLists(): Pick<
  PublicIncidentRecord,
  | "affectedServices"
  | "regions"
  | "networks"
  | "evidence"
  | "hypotheses"
  | "timeline"
  | "affectedUserCount"
  | "confidence"
> {
  return {
    affectedServices: [],
    regions: [],
    networks: [],
    evidence: [],
    hypotheses: [],
    timeline: [],
    affectedUserCount: null,
    confidence: {
      level: null,
      percent: null,
      supportingEvidenceIds: [],
      alternativeHypothesisIds: [],
      caveat: CONFIDENCE_CAVEAT,
    },
  };
}
