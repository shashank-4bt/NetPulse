import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";

import { EmptyState } from "@/components/feedback/empty-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { DiagnoseRunView } from "@/features/diagnose/diagnose-run-view";
import { SupportReportView } from "@/features/intelligence/support-report-view";
import { parseReportId } from "@/lib/reports/id";
import { loadDiagnosis } from "@/lib/reports/load";

export const dynamic = "force-dynamic";

type ReportPageProps = {
  params: Promise<{ id: string }>;
};

export async function generateMetadata({
  params,
}: ReportPageProps): Promise<Metadata> {
  const { id } = await params;
  const reportId = parseReportId(id);
  const report = reportId ? (await loadDiagnosis(reportId)).report : null;

  return {
    title: report ? `Report · ${report.target.hostname}` : "Report",
    description:
      "Diagnostic report. Ephemeral until a persistent store exists. Not indexed.",
    robots: { index: false, follow: false },
  };
}

export default async function ReportPage({ params }: ReportPageProps) {
  const { id } = await params;
  const reportId = parseReportId(id);
  if (!reportId) {
    notFound();
  }

  const loaded = await loadDiagnosis(reportId);
  const report = loaded.report;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Report"
        title={report ? report.target.hostname : "Report unavailable"}
        description={
          report
            ? "Shareable diagnostic report for web, support, future PDF, and a future JSON API. Empty fields mean the engine did not produce those values."
            : loaded.pending
              ? "The diagnosis is queued or running. Refresh when workers finish. No substitute result is shown while it is pending."
              : loaded.backendError
                ? loaded.backendError
                : "This id is well-formed, but no report is stored."
        }
      />
      <PageContainer className="space-y-10 py-10">
        <DevelopmentBanner
          title="Report store"
          description="If NETPULSE_API_BASE_URL is set, this page reads GET /v1/diagnoses/:id. Otherwise it uses the process-local Next.js store."
        />
        {report ? (
          <Tabs defaultValue="web">
            <TabsList>
              <TabsTrigger value="web">Web</TabsTrigger>
              <TabsTrigger value="support">Support technician</TabsTrigger>
            </TabsList>
            <TabsContent value="web" className="mt-6">
              <DiagnoseRunView report={report} showReportLink={false} />
            </TabsContent>
            <TabsContent value="support" className="mt-6">
              <SupportReportView report={report} />
            </TabsContent>
          </Tabs>
        ) : loaded.pending ? (
          <EmptyState
            title="Diagnosis in progress"
            description="Workers are still measuring from a NetPulse vantage. This page does not invent a result while the job is queued or running."
          />
        ) : loaded.backendError ? (
          <EmptyState
            title="Report unavailable"
            description={loaded.backendError}
          >
            <Button nativeButton={false} render={<Link href="/diagnose" />}>
              Start a diagnosis
            </Button>
          </EmptyState>
        ) : (
          <EmptyState
            title="Report not found"
            description="The store does not contain this report. That is not a measurement failure."
          >
            <Button nativeButton={false} render={<Link href="/diagnose" />}>
              Start a diagnosis
            </Button>
          </EmptyState>
        )}
      </PageContainer>
    </main>
  );
}
