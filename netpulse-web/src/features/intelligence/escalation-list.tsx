import { EmptyState } from "@/components/feedback/empty-state";
import { Badge } from "@/components/ui/badge";
import type { EscalationCondition } from "@/domain/diagnostic";

type EscalationListProps = {
  conditions: EscalationCondition[];
};

export function EscalationList({ conditions }: EscalationListProps) {
  return (
    <section aria-labelledby="escalation-heading" className="space-y-3">
      <h2 id="escalation-heading" className="text-lg font-semibold">
        Escalation conditions
      </h2>
      {conditions.length === 0 ? (
        <EmptyState
          title="No escalation conditions"
          description="Escalation rules appear only when the engine defines them. Nothing is dispatched automatically."
        />
      ) : (
        <ul className="grid gap-3">
          {conditions.map((item) => (
            <li
              key={item.id}
              className="space-y-2 rounded-lg border border-border bg-card p-4 text-sm"
            >
              <Badge variant="outline">
                {item.safetyClass === "dangerous"
                  ? "Dangerous — not executed"
                  : "Advisory"}
              </Badge>
              <p>
                <span className="font-medium">When: </span>
                <span className="text-muted-foreground">{item.when}</span>
              </p>
              <p>
                <span className="font-medium">Action: </span>
                <span className="text-muted-foreground">{item.action}</span>
              </p>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
