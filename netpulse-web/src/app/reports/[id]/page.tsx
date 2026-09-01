import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";

import { EmptyState } from "@/components/feedback/empty-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Button } from "@/components/ui/button";
import { DiagnoseRunView } from "@/features/diagnose/diagnose-run-view";
import { parseReportId } from "@/lib/reports/id";
import { getReport } from "@/lib/reports/store";

export const dynamic = "force-dynamic";

type ReportPageProps = {
  params: Promise<{ id: string }>;
};

export async function generateMetadata({
  params,
}: ReportPageProps): Promise<Metadata> {
  const { id } = await params;
  const reportId = parseReportId(id);
  const report = reportId ? getReport(reportId) : null;

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

  const report = getReport(reportId);

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Report"
        title={report ? report.target.hostname : "Report unavailable"}
        description={
          report
            ? "This page shows the stored report model. Empty fields mean the engine did not produce those values."
            : "This id is well-formed, but the report is not in this process."
        }
      />
      <PageContainer className="space-y-10 py-10">
        <DevelopmentBanner
          title="Process-local store"
          description="Reports live in memory on this Node process. They disappear on restart and are not shared across instances. PostgreSQL is not connected."
        />
        {report ? (
          <DiagnoseRunView report={report} showReportLink={false} />
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
