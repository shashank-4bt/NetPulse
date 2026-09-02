import type { Metadata } from "next";
import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { DeveloperUnavailable } from "@/features/developer/unavailable";
import { requireAccount } from "@/lib/account/guard";
import { getDevIncidents } from "@/lib/api/developer";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Developer incidents",
  robots: { index: false, follow: false },
};

export default async function DeveloperIncidentsPage() {
  const { unavailable } = await requireAccount("/developers/incidents");
  if (unavailable) {
    return <DeveloperUnavailable title="Incidents" description="Tenant monitor incidents need the API." />;
  }
  const token = await readSessionToken();
  const loaded = token ? await getDevIncidents(token) : null;
  const incidents = loaded && loaded.ok ? loaded.incidents : [];

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Developers"
        title="Monitor incidents"
        description="These rows are tenant-owned monitor incidents from stored checks, not the public catalog."
      />
      <PageContainer className="py-10">
        {!loaded || !loaded.ok ? (
          <DevelopmentBanner description={loaded && !loaded.ok ? loaded.message : "Incidents unavailable."} />
        ) : incidents.length === 0 ? (
          <p className="text-sm text-muted-foreground">No monitor incidents are stored.</p>
        ) : (
          <ul className="space-y-3 text-sm">
            {incidents.map((item) => (
              <li key={item.id}>
                <p className="font-medium">{item.title}</p>
                <p className="text-muted-foreground">
                  {item.status} · {item.startedAt} · samples {item.sampleCount}
                </p>
                <p className="text-muted-foreground">{item.summary}</p>
                <Link href={`/developers/monitors/${item.monitorId}`} className="hover:underline">
                  Open monitor
                </Link>
              </li>
            ))}
          </ul>
        )}
      </PageContainer>
    </main>
  );
}
