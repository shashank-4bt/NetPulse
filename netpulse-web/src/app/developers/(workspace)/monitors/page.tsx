import type { Metadata } from "next";
import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CreateMonitorForm } from "@/features/developer/developer-actions";
import { DeveloperUnavailable } from "@/features/developer/unavailable";
import { requireAccount } from "@/lib/account/guard";
import { getDevMonitors } from "@/lib/api/developer";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Monitors",
  robots: { index: false, follow: false },
};

export default async function DeveloperMonitorsPage() {
  const { unavailable } = await requireAccount("/developers/monitors");
  if (unavailable) {
    return <DeveloperUnavailable title="Monitors" description="HTTP, DNS, and TLS monitors need the API." />;
  }
  const token = await readSessionToken();
  const loaded = token ? await getDevMonitors(token) : null;
  const monitors = loaded && loaded.ok ? loaded.monitors : [];

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Developers"
        title="Monitors"
        description="HTTP, DNS, and TLS targets. Regions are requested vantages, not proof those regions ran."
      />
      <PageContainer className="grid gap-6 py-10 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Stored monitors</CardTitle>
          </CardHeader>
          <CardContent>
            {!loaded || !loaded.ok ? (
              <DevelopmentBanner description={loaded && !loaded.ok ? loaded.message : "Monitors unavailable."} />
            ) : monitors.length === 0 ? (
              <p className="text-sm text-muted-foreground">No monitors are stored.</p>
            ) : (
              <ul className="space-y-3 text-sm">
                {monitors.map((item) => (
                  <li key={item.id}>
                    <Link href={`/developers/monitors/${item.id}`} className="font-medium hover:underline">
                      {item.name}
                    </Link>
                    <p className="text-muted-foreground">
                      {item.type} · {item.target} · {item.status} · {item.checkCount} checks
                    </p>
                    <p className="text-muted-foreground">{item.summary}</p>
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Create monitor</CardTitle>
          </CardHeader>
          <CardContent>
            <CreateMonitorForm />
          </CardContent>
        </Card>
      </PageContainer>
    </main>
  );
}
