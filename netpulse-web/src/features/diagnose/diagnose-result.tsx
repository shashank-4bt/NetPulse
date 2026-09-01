import Link from "next/link";

import { EvidenceItem } from "@/components/data/evidence-item";
import { EmptyState } from "@/components/feedback/empty-state";
import { InsufficientEvidenceState } from "@/components/feedback/insufficient-evidence-state";
import { UnavailableState } from "@/components/feedback/unavailable-state";
import { ConfidenceBadge } from "@/components/status/confidence-badge";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import type { DiagnosticOutcome, DiagnosticReport } from "@/domain/diagnostic";

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

      <div className="grid gap-4 lg:grid-cols-2">
        <article className="rounded-lg border border-border bg-card p-4">
          <h3 className="text-sm font-medium">Likely cause</h3>
          {report.likelyCause ? (
            <p className="mt-2 text-sm text-muted-foreground">
              {report.likelyCause}
            </p>
          ) : (
            <div className="mt-3">
              <InsufficientEvidenceState description="No measured observations exist, so NetPulse will not name a cause." />
            </div>
          )}
        </article>
        <article className="rounded-lg border border-border bg-card p-4">
          <h3 className="text-sm font-medium">Confidence</h3>
          {report.confidence.level ? (
            <div className="mt-3">
              <ConfidenceBadge confidence={report.confidence.level} />
              {report.confidence.percent !== null ? (
                <p className="mt-2 font-mono text-sm tabular-nums">
                  {report.confidence.percent}%
                </p>
              ) : null}
            </div>
          ) : (
            <p className="mt-2 text-sm text-muted-foreground" role="status">
              Confidence is not available. A percentage is shown only when the
              engine supplies one.
            </p>
          )}
        </article>
      </div>

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
        {report.hypotheses.length === 0 ? (
          <EmptyState
            className="mt-3"
            title="No hypotheses"
            description="Alternatives are listed only when the engine infers more than one cause from evidence."
          />
        ) : (
          <ul className="mt-3 grid gap-3">
            {report.hypotheses.map((item) => (
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

      <div className="grid gap-4 md:grid-cols-2">
        <article className="rounded-lg border border-border bg-card p-4">
          <h3 className="text-sm font-medium">Recommended action</h3>
          <p className="mt-2 text-sm text-muted-foreground">
            {report.recommendation ?? "No recommendation without evidence."}
          </p>
        </article>
        <article className="rounded-lg border border-border bg-card p-4">
          <h3 className="text-sm font-medium">Verification</h3>
          <p className="mt-2 text-sm text-muted-foreground">
            {report.verification.note}
          </p>
        </article>
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
