import { ChartContainer } from "@/components/data/chart-container";
import { Badge } from "@/components/ui/badge";
import type { DiagnosticReport } from "@/domain/diagnostic";

type TechnicalDetailsProps = {
  report: DiagnosticReport;
};

export function TechnicalDetails({ report }: TechnicalDetailsProps) {
  return (
    <section aria-labelledby="tech-heading" className="space-y-3">
      <h2 id="tech-heading" className="text-lg font-semibold">
        Technical view
      </h2>
      <p className="text-sm text-muted-foreground">
        DNS, TCP, TLS, HTTP, latency, packet loss, route metadata, network/ASN,
        region, and timestamp. Expand a block only to inspect recorded values.
      </p>
      <div className="space-y-2">
        {report.measurements.map((block) => (
          <details
            key={block.id}
            className="rounded-lg border border-border bg-card p-3"
          >
            <summary className="cursor-pointer text-sm font-medium">
              {block.label}
              <span className="ml-2 text-xs font-normal text-muted-foreground">
                {block.measured ? "Measured" : "Not measured"}
              </span>
            </summary>
            <dl className="mt-3 grid gap-2 text-sm text-muted-foreground">
              <div className="flex flex-wrap items-center gap-2">
                <dt className="sr-only">State</dt>
                <dd>
                  <Badge variant="outline">
                    {block.measured ? "Measured" : "Not measured"}
                  </Badge>
                </dd>
              </div>
              <div>
                <dt className="font-medium text-foreground">Value</dt>
                <dd className="font-mono text-xs">
                  {block.value === null
                    ? "No value recorded"
                    : `${block.value}${block.unit ? ` ${block.unit}` : ""}`}
                </dd>
              </div>
              <div>
                <dt className="font-medium text-foreground">Measured at</dt>
                <dd className="font-mono text-xs">
                  {block.measuredAt ?? "No measurement timestamp"}
                </dd>
              </div>
              {block.summary ? (
                <div>
                  <dt className="font-medium text-foreground">Summary</dt>
                  <dd className="font-mono text-xs">{block.summary}</dd>
                </div>
              ) : null}
            </dl>
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
