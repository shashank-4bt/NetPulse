import type { Metadata } from "next";

import { EmptyState } from "@/components/feedback/empty-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { LoadingState } from "@/components/feedback/loading-state";
import { UnavailableState } from "@/components/feedback/unavailable-state";
import { DiagnoseForm } from "@/features/diagnose/diagnose-form";
import { DiagnoseRunView } from "@/features/diagnose/diagnose-run-view";
import { DiagnoseWorkflow } from "@/features/diagnose/diagnose-workflow";
import { safeDiagnosePrefill } from "@/features/diagnose/safe-prefill";
import { isApiConfigured } from "@/lib/api/backend";
import { parseReportId } from "@/lib/reports/id";
import { loadDiagnosis } from "@/lib/reports/load";

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
  const loaded = reportId ? await loadDiagnosis(reportId) : null;
  const report = loaded?.report ?? null;
  const missingRun = Boolean(params.run) && Boolean(loaded?.missing);
  const initialValue = safeDiagnosePrefill(params.target);
  const apiConfigured = isApiConfigured();

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Diagnosis"
        title="Check a public target"
        description="Enter a domain, URL, or known service. NetPulse will not invent a cause, confidence, or percentage when probes have not run."
      />
      <PageContainer className="space-y-10 py-10">
        <DevelopmentBanner
          title={
            apiConfigured
              ? "Worker vantage measurements"
              : "Go API is not configured"
          }
          description={
            apiConfigured
              ? "When the API is up, workers record DNS, TCP, TLS, and HTTP from NetPulse — not from your device. That cannot isolate a user-path root cause by itself."
              : "Set NETPULSE_API_BASE_URL to connect /diagnose to POST /v1/diagnoses. Until then, submits stay local and measurement-unavailable."
          }
        />
        <DiagnoseForm initialValue={initialValue} />
        {missingRun ? (
          <EmptyState
            title="Diagnosis not found"
            description="The API or process store does not contain this id. That is not a measurement failure."
          />
        ) : null}
        {loaded?.backendError ? (
          <UnavailableState
            title="Backend unavailable"
            description={loaded.backendError}
          />
        ) : null}
        {loaded?.pending ? (
          <LoadingState label="Waiting for measurement workers" />
        ) : null}
        {report ? <DiagnoseRunView report={report} /> : <DiagnoseWorkflow />}
      </PageContainer>
    </main>
  );
}
