import Link from "next/link";

import { EvidenceItem } from "@/components/data/evidence-item";
import { EmptyState } from "@/components/feedback/empty-state";
import { UnavailableState } from "@/components/feedback/unavailable-state";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import type { DiagnosticOutcome, DiagnosticReport } from "@/domain/diagnostic";
import { InsufficientEvidenceResultView } from "@/features/intelligence/insufficient-evidence-result";

const OUTCOME_LABEL: Record<DiagnosticOutcome, string> = {
  success: "Success",
  partial_success: "Partial success",
  timeout: "Timeout",
  backend_unavailable: "Backend unavailable",
  insufficient_evidence: "Insufficient evidence",
  invalid_input: "Invalid input",
  measurement_unavailable: "Measurement unavailable",
};

type DiagnoseResultProps = {
  report: DiagnosticReport;
  showReportLink?: boolean;
};

export function DiagnoseResult({
  report,
  showReportLink = true,
}: DiagnoseResultProps) {
  return (
    <section aria-labelledby="result-heading" className="space-y-6">
      <div className="flex flex-wrap items-center gap-2">
        <h2 id="result-heading" className="text-lg font-semibold">
          Result
        </h2>
        <Badge variant="outline">{OUTCOME_LABEL[report.outcome]}</Badge>
      </div>

      <InsufficientEvidenceResultView result={report.insufficientEvidence} />

      <article className="rounded-lg border border-border bg-card p-4">
        <h3 className="text-sm font-medium">Likely cause</h3>
        {report.likelyCause ? (
          <p className="mt-2 text-sm text-muted-foreground">
            {report.likelyCause}
          </p>
        ) : (
          <p className="mt-2 text-sm text-muted-foreground" role="status">
            No cause is named. That is insufficient evidence, not a fabricated
            root cause.
          </p>
        )}
      </article>

      <div>
        <h3 className="text-sm font-medium">Evidence</h3>
        {report.evidence.length === 0 ? (
          <EmptyState
            className="mt-3"
            title="No evidence recorded"
            description="Facts appear here only after probes persist observations."
          />
        ) : (
          <ul className="mt-3 grid gap-3 lg:grid-cols-2">
            {report.evidence.map((item) => (
              <li key={item.id}>
                <EvidenceItem
                  evidenceClass={item.evidenceClass}
                  title={item.title}
                  body={item.body}
                />
              </li>
            ))}
          </ul>
        )}
      </div>

      <div>
        <h3 className="text-sm font-medium">Alternative hypotheses</h3>
        {report.alternativeHypotheses.length === 0 ? (
          <EmptyState
            className="mt-3"
            title="No hypotheses"
            description="Alternatives are listed only when the engine infers more than one cause from evidence."
          />
        ) : (
          <ul className="mt-3 grid gap-3">
            {report.alternativeHypotheses.map((item) => (
              <li key={item.id}>
                <EvidenceItem
                  evidenceClass="inferred_hypothesis"
                  title={item.title}
                  body={item.body}
                />
              </li>
            ))}
          </ul>
        )}
      </div>

      <div>
        <h3 className="text-sm font-medium">Verification steps</h3>
        {report.verificationSteps.length === 0 ? (
          <EmptyState
            className="mt-3"
            title="No verification steps"
            description="Verification compares a later measured run to this one."
          />
        ) : (
          <ul className="mt-3 grid gap-3">
            {report.verificationSteps.map((step) => (
              <li
                key={step.id}
                className="rounded-lg border border-border bg-card p-4 text-sm"
              >
                <div className="flex flex-wrap items-center gap-2">
                  <p className="font-medium">{step.label}</p>
                  <Badge variant="outline">{step.status.replace(/_/g, " ")}</Badge>
                </div>
                <p className="mt-2 text-muted-foreground">{step.note}</p>
              </li>
            ))}
          </ul>
        )}
      </div>

      {report.outcome === "measurement_unavailable" ||
      report.outcome === "backend_unavailable" ? (
        <UnavailableState
          title="Workers are not connected"
          description="This result documents an accepted target and an unavailable measurement path. It is not a live internet diagnosis."
        />
      ) : null}

      {showReportLink ? (
        <Button
          variant="outline"
          nativeButton={false}
          render={<Link href={`/reports/${report.reportId}`} />}
        >
          Open report
        </Button>
      ) : null}
    </section>
  );
}
