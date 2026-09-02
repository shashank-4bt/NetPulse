import type { Metadata } from "next";
import Link from "next/link";

import { EmptyState } from "@/components/feedback/empty-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { HistoryFilters } from "@/features/account/history-filters";
import { parseHistoryQuery } from "@/features/account/history";
import { requireAccount } from "@/lib/account/guard";
import { getHistory } from "@/lib/api/account";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Diagnosis history",
  robots: { index: false, follow: false },
};

type PageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export default async function HistoryPage({ searchParams }: PageProps) {
  const { unavailable } = await requireAccount("/dashboard/history");
  const query = parseHistoryQuery(await searchParams);
  if (unavailable) {
    return (
      <main id="main-content" className="flex-1">
        <PageHero eyebrow="History" title="Diagnosis history" description="Search and filters need the API." />
        <PageContainer className="py-10">
          <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set. History is empty." />
        </PageContainer>
      </main>
    );
  }
  const token = await readSessionToken();
  const loaded = token ? await getHistory(token, query) : null;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="History"
        title="Diagnosis history"
        description="Search, status, target, and date range apply only to diagnoses you started while signed in."
      />
      <PageContainer className="space-y-6 py-10">
        <HistoryFilters query={query} />
        {!loaded || !loaded.ok ? (
          <EmptyState
            title="History unavailable"
            description={loaded && !loaded.ok ? loaded.message : "Sign in required."}
          />
        ) : loaded.diagnoses.length === 0 ? (
          <EmptyState
            title="No matching diagnoses"
            description="Nothing stored matches these filters. That is not a live outage."
          />
        ) : (
          <ul className="divide-y divide-border rounded-lg border border-border">
            {loaded.diagnoses.map((item) => (
              <li key={item.id} className="px-4 py-3 text-sm">
                <Link href={`/reports/${item.id}`} className="font-medium hover:underline">
                  {item.target}
                </Link>
                <p className="text-muted-foreground">
                  {item.status}
                  {item.outcome ? ` · ${item.outcome}` : ""} · {item.createdAt}
                </p>
              </li>
            ))}
          </ul>
        )}
      </PageContainer>
    </main>
  );
}
