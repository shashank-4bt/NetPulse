import { DiagnoseResult } from "@/features/diagnose/diagnose-result";
import { DiagnoseStepper } from "@/features/diagnose/diagnose-stepper";
import { TechnicalDetails } from "@/features/diagnose/technical-details";
import { IntelligenceView } from "@/features/intelligence/intelligence-view";
import type { DiagnosticReport } from "@/domain/diagnostic";

type DiagnoseRunViewProps = {
  report: DiagnosticReport;
  showReportLink?: boolean;
};

export function DiagnoseRunView({
  report,
  showReportLink = true,
}: DiagnoseRunViewProps) {
  return (
    <div className="space-y-10">
      <section aria-labelledby="run-meta-heading" className="space-y-2">
        <h2 id="run-meta-heading" className="text-lg font-semibold">
          Run
        </h2>
        <dl className="grid gap-2 text-sm sm:grid-cols-2">
          <div>
            <dt className="text-muted-foreground">Target</dt>
            <dd className="font-mono">{report.target.hostname}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">Input kind</dt>
            <dd>{report.target.kind.replace(/_/g, " ")}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">Timestamp</dt>
            <dd className="font-mono">{report.timestamp}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">diagnosticEngineVersion</dt>
            <dd className="font-mono">{report.versions.diagnosticEngineVersion}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">ruleVersion</dt>
            <dd className="font-mono">{report.versions.ruleVersion}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">measurementVersion</dt>
            <dd className="font-mono">{report.versions.measurementVersion}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">modelVersion</dt>
            <dd className="font-mono">{report.versions.modelVersion}</dd>
          </div>
          <div className="sm:col-span-2">
            <dt className="text-muted-foreground">Report ID</dt>
            <dd className="font-mono break-all">{report.reportId}</dd>
          </div>
        </dl>
      </section>
      <DiagnoseStepper steps={report.tests} />
      <DiagnoseResult report={report} showReportLink={showReportLink} />
      <IntelligenceView report={report} />
      <TechnicalDetails report={report} />
    </div>
  );
}
