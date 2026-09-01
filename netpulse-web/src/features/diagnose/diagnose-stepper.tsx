import { Badge } from "@/components/ui/badge";
import type { DiagnosticStep } from "@/domain/diagnostic";
import { cn } from "@/lib/utils";

const STATE_LABEL: Record<DiagnosticStep["state"], string> = {
  pending: "Pending",
  current: "Current",
  complete: "Complete",
  failed: "Failed",
  unavailable: "Unavailable",
  not_measured: "Not measured",
};

type DiagnoseStepperProps = {
  steps: DiagnosticStep[];
};

export function DiagnoseStepper({ steps }: DiagnoseStepperProps) {
  const completeCount = steps.filter((step) => step.state === "complete").length;
  const unavailableCount = steps.filter(
    (step) => step.state === "unavailable" || step.state === "not_measured"
  ).length;
  const failedCount = steps.filter((step) => step.state === "failed").length;

  return (
    <section aria-labelledby="stepper-heading" className="space-y-4">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <h2 id="stepper-heading" className="text-lg font-semibold">
          Run steps
        </h2>
        <p className="text-sm text-muted-foreground">
          {completeCount} complete · {unavailableCount} unavailable
          {failedCount > 0 ? ` · ${failedCount} failed` : ""} · {steps.length}{" "}
          total
        </p>
      </div>
      <ol className="grid gap-2 sm:grid-cols-2">
        {steps.map((step, index) => (
          <li
            key={step.id}
            className={cn(
              "flex flex-col gap-1 rounded-lg border border-border bg-card p-3",
              step.state === "current" && "ring-2 ring-ring"
            )}
          >
            <div className="flex items-center justify-between gap-2">
              <p className="font-mono text-xs text-muted-foreground">
                {String(index + 1).padStart(2, "0")}
              </p>
              <Badge variant="outline">{STATE_LABEL[step.state]}</Badge>
            </div>
            <p className="text-sm font-medium">{step.label}</p>
            {step.durationMs !== null ? (
              <p className="font-mono text-xs text-muted-foreground">
                {step.durationMs} ms
              </p>
            ) : null}
            {step.note ? (
              <p className="text-xs text-muted-foreground">{step.note}</p>
            ) : null}
          </li>
        ))}
      </ol>
    </section>
  );
}
