import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { IncidentList, LatencyCard, ObservationCard, RegionalList } from "@/features/developer/metric-cards";
import { DeveloperUnavailable } from "@/features/developer/unavailable";
import { requireAccount } from "@/lib/account/guard";
import { getDevSLA } from "@/lib/api/developer";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "SLA",
  robots: { index: false, follow: false },
};

export default async function DeveloperSLAPage() {
  const { unavailable } = await requireAccount("/developers/sla");
  if (unavailable) {
    return <DeveloperUnavailable title="SLA" description="SLA ratios need stored monitor checks." />;
  }
  const token = await readSessionToken();
  const loaded = token ? await getDevSLA(token) : null;
  const sla = loaded && loaded.ok ? loaded.sla : null;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Developers"
        title="SLA"
        description="Failed-versus-stored-check ratios. This is not calendar uptime and not a billed guarantee."
      />
      <PageContainer className="space-y-6 py-10">
        {!sla ? (
          <DevelopmentBanner description={loaded && !loaded.ok ? loaded.message : "SLA unavailable."} />
        ) : (
          <>
            <p className="text-sm text-muted-foreground">{sla.summary}</p>
            <div className="grid gap-4 md:grid-cols-2">
              <ObservationCard title="Availability" item={sla.availability} />
              <ObservationCard title="Downtime" item={sla.downtime} />
              <LatencyCard item={sla.latency} />
              <Card>
                <CardHeader>
                  <CardTitle>Incidents</CardTitle>
                </CardHeader>
                <CardContent>
                  <IncidentList items={sla.incidents} />
                </CardContent>
              </Card>
              <Card className="md:col-span-2">
                <CardHeader>
                  <CardTitle>Regional performance</CardTitle>
                </CardHeader>
                <CardContent>
                  <RegionalList items={sla.regionalPerformance} />
                </CardContent>
              </Card>
            </div>
          </>
        )}
      </PageContainer>
    </main>
  );
}
