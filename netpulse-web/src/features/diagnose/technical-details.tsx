import { ChartContainer } from "@/components/data/chart-container";
import type { DiagnosticReport } from "@/domain/diagnostic";

type TechnicalDetailsProps = {
  report: DiagnosticReport;
};

export function TechnicalDetails({ report }: TechnicalDetailsProps) {
  return (
    <section aria-labelledby="tech-heading" className="space-y-3">
      <h2 id="tech-heading" className="text-lg font-semibold">
        Technical details
      </h2>
      <p className="text-sm text-muted-foreground">
        Expand a block only to inspect recorded values. Empty blocks mean the
        probe did not run.
      </p>
      <div className="space-y-2">
        {report.measurements.map((block) => (
          <details
            key={block.key}
            className="rounded-lg border border-border bg-card p-3"
          >
            <summary className="cursor-pointer text-sm font-medium">
              {block.label}
              <span className="ml-2 text-xs font-normal text-muted-foreground">
                {block.measured ? "Measured" : "Not measured"}
              </span>
            </summary>
            <div className="mt-3 text-sm text-muted-foreground">
              {block.summary ? (
                <p className="font-mono text-xs">{block.summary}</p>
              ) : (
                <p>No values were recorded for {block.label}.</p>
              )}
            </div>
          </details>
        ))}
      </div>
      <ChartContainer
        title="Latency series"
        description="A chart is rendered only when a measured series exists."
        state="unavailable"
      />
    </section>
  );
}
