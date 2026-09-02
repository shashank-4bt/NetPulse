import type { Metadata } from "next";
import Link from "next/link";

import { EmptyState } from "@/components/feedback/empty-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import {
  DeleteReportButton,
  ExportReportButton,
  ShareReportButton,
} from "@/features/account/account-actions";
import { requireAccount } from "@/lib/account/guard";
import { getReports } from "@/lib/api/account";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Reports",
  robots: { index: false, follow: false },
};

export default async function ReportsPage() {
  const { unavailable } = await requireAccount("/dashboard/reports");
  if (unavailable) {
    return (
      <main id="main-content" className="flex-1">
        <PageHero eyebrow="Reports" title="Your reports" description="Report actions need the API." />
        <PageContainer className="py-10">
          <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set. No sample reports are shown." />
        </PageContainer>
      </main>
    );
  }
  const token = await readSessionToken();
  const loaded = token ? await getReports(token) : null;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Reports"
        title="Your reports"
        description="View, share, delete, or download JSON. PDF export is not ready."
      />
      <PageContainer className="space-y-6 py-10">
        {!loaded || !loaded.ok ? (
          <EmptyState
            title="Reports unavailable"
            description={loaded && !loaded.ok ? loaded.message : "Sign in required."}
          />
        ) : loaded.reports.length === 0 ? (
          <EmptyState title="No reports" description="Diagnoses you start while signed in appear here." />
        ) : (
          <ul className="space-y-4">
            {loaded.reports.map((item) => (
              <li key={item.id} className="space-y-3 rounded-lg border border-border p-4">
                <div>
                  <Link href={`/reports/${item.id}`} className="font-medium hover:underline">
                    {item.target}
                  </Link>
                  <p className="text-sm text-muted-foreground">
                    {item.status}
                    {item.shared ? " · share link exists" : ""} · {item.createdAt}
                  </p>
                </div>
                <div className="flex flex-wrap gap-3">
                  <ShareReportButton reportId={item.id} />
                  <ExportReportButton reportId={item.id} />
                  <DeleteReportButton reportId={item.id} />
                </div>
              </li>
            ))}
          </ul>
        )}
      </PageContainer>
    </main>
  );
}
