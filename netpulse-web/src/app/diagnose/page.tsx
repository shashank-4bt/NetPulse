import type { Metadata } from "next";

import { EmptyState } from "@/components/feedback/empty-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { DiagnoseForm } from "@/features/diagnose/diagnose-form";
import { DiagnoseRunView } from "@/features/diagnose/diagnose-run-view";
import { DiagnoseWorkflow } from "@/features/diagnose/diagnose-workflow";
import { safeDiagnosePrefill } from "@/features/diagnose/safe-prefill";
import { parseReportId } from "@/lib/reports/id";
import { getReport } from "@/lib/reports/store";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Diagnose",
  description:
    "Start a NetPulse diagnosis for a public hostname, URL, or known service. Measurements stay unavailable until workers are connected.",
  alternates: { canonical: "/diagnose" },
};

type DiagnosePageProps = {
  searchParams: Promise<{ run?: string; target?: string }>;
};

export default async function DiagnosePage({ searchParams }: DiagnosePageProps) {
  const params = await searchParams;
  const reportId = parseReportId(params.run);
  const report = reportId ? getReport(reportId) : null;
  const missingRun = Boolean(params.run) && !report;
  const initialValue = safeDiagnosePrefill(params.target);

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Diagnosis"
        title="Check a public target"
        description="Enter a domain, URL, or known service. NetPulse will not invent a cause, confidence, or percentage when probes have not run."
      />
      <PageContainer className="space-y-10 py-10">
        <DevelopmentBanner
          title="Measurement workers are not connected"
          description="Submitting a target creates a report that records the accepted hostname and marks every probe step unavailable. That is not a live diagnosis."
        />
        <DiagnoseForm initialValue={initialValue} />
        {missingRun ? (
          <EmptyState
            title="Report not in this process"
            description="The run id is missing from the process-local store. Reports are not persisted and are not shared across instances."
          />
        ) : null}
        {report ? <DiagnoseRunView report={report} /> : <DiagnoseWorkflow />}
      </PageContainer>
    </main>
  );
}
