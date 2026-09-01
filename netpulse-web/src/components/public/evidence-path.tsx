import { EVIDENCE_FLOW } from "@/lib/content/evidence-flow";

export function EvidencePath() {
  return (
    <ol className="grid gap-3 md:grid-cols-4">
      {EVIDENCE_FLOW.map((step, index) => (
        <li
          key={step.id}
          className="flex flex-col gap-2 rounded-lg border border-border bg-card p-4"
        >
          <p className="font-mono text-xs text-muted-foreground">
            {String(index + 1).padStart(2, "0")}
          </p>
          <h3 className="text-sm font-medium">{step.name}</h3>
          <p className="text-sm text-muted-foreground">{step.summary}</p>
        </li>
      ))}
    </ol>
  );
}
