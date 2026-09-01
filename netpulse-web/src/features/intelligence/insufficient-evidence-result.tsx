import { InsufficientEvidenceState } from "@/components/feedback/insufficient-evidence-state";
import type { InsufficientEvidenceResult } from "@/domain/diagnostic";

type InsufficientEvidenceResultProps = {
  result: InsufficientEvidenceResult;
};

export function InsufficientEvidenceResultView({
  result,
}: InsufficientEvidenceResultProps) {
  if (!result.determined) {
    return null;
  }

  return (
    <section
      aria-labelledby="insufficient-heading"
      className="space-y-3 rounded-lg border border-border bg-card p-4"
    >
      <h2 id="insufficient-heading" className="text-lg font-semibold">
        Insufficient evidence
      </h2>
      <InsufficientEvidenceState description={result.message} />
      <p className="text-sm text-foreground">{result.nextCheck}</p>
    </section>
  );
}
