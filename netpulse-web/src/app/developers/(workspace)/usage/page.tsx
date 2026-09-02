import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DeveloperUnavailable } from "@/features/developer/unavailable";
import { requireAccount } from "@/lib/account/guard";
import { getDevUsage } from "@/lib/api/developer";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Usage",
  robots: { index: false, follow: false },
};

export default async function DeveloperUsagePage() {
  const { unavailable } = await requireAccount("/developers/usage");
  if (unavailable) {
    return <DeveloperUnavailable title="Usage" description="Usage counters need the API." />;
  }
  const token = await readSessionToken();
  const loaded = token ? await getDevUsage(token) : null;
  const usage = loaded && loaded.ok ? loaded.usage : null;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Developers"
        title="Usage"
        description="Stored request and resource counts. Traffic is not estimated."
      />
      <PageContainer className="py-10">
        {!usage ? (
          <DevelopmentBanner description={loaded && !loaded.ok ? loaded.message : "Usage unavailable."} />
        ) : (
          <div className="grid gap-4 md:grid-cols-2">
            <Stat title="API key requests" value={usage.requests} />
            <Stat title="Measurements" value={usage.measurements} />
            <Stat title="Monitors" value={usage.monitors} />
            <Stat title="Requested regions" value={usage.regions} />
            <Stat title="Webhooks" value={usage.webhooks} />
            <Card className="md:col-span-2">
              <CardHeader>
                <CardTitle>Notes</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">{usage.summary}</CardContent>
            </Card>
          </div>
        )}
      </PageContainer>
    </main>
  );
}

function Stat({ title, value }: { title: string; value: number }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent className="text-2xl font-semibold">{value}</CardContent>
    </Card>
  );
}
