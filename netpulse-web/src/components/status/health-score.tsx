import type { ConfidenceLevel } from "@/domain/display";
import { ConfidenceBadge } from "@/components/status/confidence-badge";
import { INSUFFICIENT_EVIDENCE_LABEL } from "@/lib/design/taxonomy";
import { cn } from "@/lib/utils";

type HealthScoreProps = {
  /** When null, the component must not invent a score. */
  value: number | null;
  confidence: ConfidenceLevel | null;
  label?: string;
  className?: string;
};

export function HealthScore({
  value,
  confidence,
  label = "Health score",
  className,
}: HealthScoreProps) {
  const hasScore = value !== null;

  return (
    <div
      className={cn(
        "flex flex-col gap-2 rounded-lg border border-border bg-card p-4",
        className
      )}
    >
      <p className="text-xs font-medium tracking-wide text-muted-foreground uppercase">
        {label}
      </p>
      {hasScore ? (
        <>
          <p className="font-mono text-3xl font-semibold tabular-nums text-foreground">
            {value}
            <span className="ml-1 text-base font-normal text-muted-foreground">
              / 100
            </span>
          </p>
          {confidence ? (
            <ConfidenceBadge confidence={confidence} />
          ) : (
            <p className="text-sm text-muted-foreground">
              Score shown without a confidence level.
            </p>
          )}
        </>
      ) : (
        <p className="text-sm text-muted-foreground" role="status">
          {INSUFFICIENT_EVIDENCE_LABEL}
        </p>
      )}
    </div>
  );
}
