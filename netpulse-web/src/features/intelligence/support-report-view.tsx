import { Badge } from "@/components/ui/badge";
import type { DiagnosticReport } from "@/domain/diagnostic";
import { stringifyShareableReport } from "@/features/intelligence/shareable-report";
import { LAYER_STATUS_LABEL } from "@/lib/design/taxonomy";

type SupportReportViewProps = {
  report: DiagnosticReport;
};

export function SupportReportView({ report }: SupportReportViewProps) {
  return (
    <article className="space-y-8" aria-labelledby="support-heading">
      <header className="space-y-2">
        <h2 id="support-heading" className="text-lg font-semibold">
          Support technician view
        </h2>
        <p className="text-sm text-muted-foreground">
          Same stored report as the web view. Empty fields are omitted facts, not
          hidden failures. Suitable for handoff, future PDF, and a future JSON
          API.
        </p>
      </header>

      <section aria-labelledby="support-identity-heading" className="space-y-2">
        <h3 id="support-identity-heading" className="text-sm font-medium">
          Identity
        </h3>
        <dl className="grid gap-2 font-mono text-xs sm:grid-cols-2">
          <div>
            <dt className="text-muted-foreground">reportId</dt>
            <dd className="break-all">{report.reportId}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">target</dt>
            <dd>{report.target.hostname}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">timestamp</dt>
            <dd>{report.timestamp}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">outcome</dt>
            <dd>{report.outcome}</dd>
          </div>
        </dl>
      </section>

      <section aria-labelledby="support-versions-heading" className="space-y-2">
        <h3 id="support-versions-heading" className="text-sm font-medium">
          Versions
        </h3>
        <dl className="grid gap-2 font-mono text-xs sm:grid-cols-2">
          <div>
            <dt className="text-muted-foreground">diagnosticEngineVersion</dt>
            <dd>{report.versions.diagnosticEngineVersion}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">ruleVersion</dt>
            <dd>{report.versions.ruleVersion}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">measurementVersion</dt>
            <dd>{report.versions.measurementVersion}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">modelVersion</dt>
            <dd>{report.versions.modelVersion}</dd>
          </div>
        </dl>
      </section>

      <section aria-labelledby="support-graph-heading" className="space-y-2">
        <h3 id="support-graph-heading" className="text-sm font-medium">
          Graph
        </h3>
        <ul className="grid gap-2">
          {report.graph.map((node) => (
            <li
              key={node.id}
              className="flex flex-wrap items-center justify-between gap-2 rounded-md border border-border px-3 py-2 text-sm"
            >
              <span>{node.label}</span>
              <Badge variant="outline">{LAYER_STATUS_LABEL[node.status]}</Badge>
            </li>
          ))}
        </ul>
      </section>

      <section aria-labelledby="support-json-heading" className="space-y-2">
        <h3 id="support-json-heading" className="text-sm font-medium">
          Shareable JSON
        </h3>
        <pre className="max-h-96 overflow-auto rounded-lg border border-border bg-muted/30 p-3 font-mono text-xs whitespace-pre-wrap">
          {stringifyShareableReport(report)}
        </pre>
      </section>
    </article>
  );
}
