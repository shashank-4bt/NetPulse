import { DIAGNOSTIC_STEP_IDS, DIAGNOSTIC_STEP_LABELS } from "@/domain/diagnostic";

export function DiagnoseWorkflow() {
  return (
    <section aria-labelledby="workflow-heading" className="space-y-3">
      <h2 id="workflow-heading" className="text-lg font-semibold">
        Diagnosis path
      </h2>
      <p className="text-sm text-muted-foreground">
        A run walks this order. Status badges appear only after a report exists.
        Unmeasured steps stay unavailable — they are not failures.
      </p>
      <ol className="flex flex-wrap gap-2">
        {DIAGNOSTIC_STEP_IDS.map((id, index) => (
          <li
            key={id}
            className="rounded-md border border-border bg-card px-2.5 py-1 font-mono text-xs"
          >
            {String(index + 1).padStart(2, "0")} {DIAGNOSTIC_STEP_LABELS[id]}
          </li>
        ))}
      </ol>
    </section>
  );
}
