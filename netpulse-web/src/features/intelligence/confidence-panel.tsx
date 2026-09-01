import { EvidenceItem } from "@/components/data/evidence-item";
import { EmptyState } from "@/components/feedback/empty-state";
import { ConfidenceBadge } from "@/components/status/confidence-badge";
import type { AlternativeHypothesis, Confidence, Evidence } from "@/domain/diagnostic";
import { formatConfidenceValue } from "@/features/intelligence/confidence";
import { lookupByIds } from "@/features/intelligence/honesty";

type ConfidencePanelProps = {
  confidence: Confidence;
  evidence: Evidence[];
  alternatives: AlternativeHypothesis[];
};

export function ConfidencePanel({
  confidence,
  evidence,
  alternatives,
}: ConfidencePanelProps) {
  const supporting = lookupByIds(evidence, confidence.supportingEvidenceIds);
  const alternativeExplanations = lookupByIds(
    alternatives,
    confidence.alternativeHypothesisIds
  );

  return (
    <section aria-labelledby="confidence-heading" className="space-y-4">
      <h2 id="confidence-heading" className="text-lg font-semibold">
        Confidence
      </h2>
      <p className="text-sm text-muted-foreground" role="note">
        {confidence.caveat}
      </p>
      <div className="grid gap-4 sm:grid-cols-2">
        <article className="rounded-lg border border-border bg-card p-4">
          <h3 className="text-sm font-medium">Confidence value</h3>
          <p className="mt-2 font-mono text-sm tabular-nums">
            {formatConfidenceValue(confidence)}
          </p>
        </article>
        <article className="rounded-lg border border-border bg-card p-4">
          <h3 className="text-sm font-medium">Confidence level</h3>
          {confidence.level ? (
            <div className="mt-2">
              <ConfidenceBadge confidence={confidence.level} />
            </div>
          ) : (
            <p className="mt-2 text-sm text-muted-foreground" role="status">
              No level was supplied. Absence of a level is not a failure.
            </p>
          )}
        </article>
      </div>
      <div>
        <h3 className="text-sm font-medium">Supporting evidence</h3>
        {supporting.length === 0 ? (
          <EmptyState
            className="mt-3"
            title="No supporting evidence"
            description="Evidence ids are listed here only when the engine attaches measured facts to this confidence."
          />
        ) : (
          <ul className="mt-3 grid gap-3 lg:grid-cols-2">
            {supporting.map((item) => (
              <li key={item.id}>
                <EvidenceItem
                  evidenceClass={item.evidenceClass}
                  title={item.title}
                  body={item.body}
                />
              </li>
            ))}
          </ul>
        )}
      </div>
      <div>
        <h3 className="text-sm font-medium">Alternative explanations</h3>
        {alternativeExplanations.length === 0 ? (
          <EmptyState
            className="mt-3"
            title="No alternative explanations"
            description="Alternatives appear only when the engine infers more than one hypothesis from evidence."
          />
        ) : (
          <ul className="mt-3 grid gap-3">
            {alternativeExplanations.map((item) => (
              <li key={item.id}>
                <EvidenceItem
                  evidenceClass="inferred_hypothesis"
                  title={item.title}
                  body={item.body}
                />
              </li>
            ))}
          </ul>
        )}
      </div>
    </section>
  );
}
