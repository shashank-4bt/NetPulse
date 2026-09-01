import type {
  DiagnosticReport,
  Evidence,
  Recommendation,
} from "@/domain/diagnostic";

export function measuredFacts(report: DiagnosticReport): Evidence[] {
  return report.evidence.filter((item) => item.evidenceClass === "measured_fact");
}

export function lookupByIds<T extends { id: string }>(
  items: readonly T[],
  ids: readonly string[]
): T[] {
  return ids
    .map((id) => items.find((item) => item.id === id))
    .filter((item): item is T => item !== undefined);
}

export function reportHonestyErrors(report: DiagnosticReport): string[] {
  const errors: string[] = [];
  const facts = measuredFacts(report);

  if (facts.length === 0) {
    if (report.likelyCause !== null) {
      errors.push("A likely cause requires measured facts.");
    }
    if (report.confidence.percent !== null || report.confidence.level !== null) {
      errors.push("Confidence requires measured facts.");
    }
    if (report.hypotheses.length > 0 || report.alternativeHypotheses.length > 0) {
      errors.push("Hypotheses require measured facts.");
    }
    if (report.graph.some((node) => node.status === "failed")) {
      errors.push("Unknown or unmeasured nodes must not be marked failed.");
    }
    if (report.graph.some((node) => node.status === "healthy" || node.status === "degraded")) {
      errors.push("Layer health requires measured facts.");
    }
    if (!report.insufficientEvidence.determined) {
      errors.push("Missing facts must produce a first-class insufficient-evidence result.");
    }
  }

  for (const node of report.graph) {
    if (node.status === "failed" && node.evidenceIds.length === 0) {
      errors.push(`Node ${node.id} is failed without evidence.`);
    }
  }

  for (const recommendation of report.recommendations) {
    if (recommendation.autoExecute) {
      errors.push(`Recommendation ${recommendation.id} must not auto-execute.`);
    }
    if (isDangerous(recommendation) && recommendation.autoExecute) {
      errors.push(`Dangerous recommendation ${recommendation.id} must never auto-execute.`);
    }
  }

  if (report.confidence.caveat.trim() === "") {
    errors.push("Confidence must include a no-certainty caveat.");
  }

  return errors;
}

export function isDangerous(recommendation: Recommendation): boolean {
  return recommendation.safetyClass === "dangerous";
}
