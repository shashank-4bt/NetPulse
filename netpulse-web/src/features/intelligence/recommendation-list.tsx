import { EmptyState } from "@/components/feedback/empty-state";
import { Badge } from "@/components/ui/badge";
import type { Recommendation } from "@/domain/diagnostic";

const SAFETY_LABEL: Record<Recommendation["safetyClass"], string> = {
  safe: "Safe",
  advisory: "Advisory",
  dangerous: "Dangerous — not executed",
};

type RecommendationListProps = {
  recommendations: Recommendation[];
};

export function RecommendationList({ recommendations }: RecommendationListProps) {
  return (
    <section aria-labelledby="recommendation-heading" className="space-y-3">
      <h2 id="recommendation-heading" className="text-lg font-semibold">
        Recommendations
      </h2>
      <p className="text-sm text-muted-foreground">
        NetPulse displays recommended actions. It never executes them, including
        dangerous ones.
      </p>
      {recommendations.length === 0 ? (
        <EmptyState
          title="No recommendations"
          description="Recommendations are listed only when the engine can attach an action to evidence."
        />
      ) : (
        <ul className="grid gap-3">
          {recommendations.map((item) => (
            <li
              key={item.id}
              className="space-y-3 rounded-lg border border-border bg-card p-4"
            >
              <div className="flex flex-wrap items-center gap-2">
                <Badge variant="outline">{SAFETY_LABEL[item.safetyClass]}</Badge>
                <Badge variant="outline">Not auto-executed</Badge>
              </div>
              {item.safetyClass === "dangerous" ? (
                <p className="text-sm text-muted-foreground">
                  Possible security concern. NetPulse does not execute this
                  action. Confirm the evidence, confidence, and a recovery plan
                  before you act. Do not delete files, change registry settings,
                  disable security tools, factory-reset, or flash firmware from
                  this report.
                </p>
              ) : null}
              <dl className="grid gap-3 text-sm">
                <div>
                  <dt className="font-medium">Action</dt>
                  <dd className="mt-1 text-muted-foreground">{item.action}</dd>
                </div>
                <div>
                  <dt className="font-medium">Reason</dt>
                  <dd className="mt-1 text-muted-foreground">{item.reason}</dd>
                </div>
                <div>
                  <dt className="font-medium">Risk</dt>
                  <dd className="mt-1 text-muted-foreground">{item.risk}</dd>
                </div>
                <div>
                  <dt className="font-medium">Expected result</dt>
                  <dd className="mt-1 text-muted-foreground">
                    {item.expectedResult}
                  </dd>
                </div>
                <div>
                  <dt className="font-medium">Verification</dt>
                  <dd className="mt-1 text-muted-foreground">
                    {item.verification}
                  </dd>
                </div>
              </dl>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
