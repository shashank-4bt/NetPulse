import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CreateAlertForm, DeleteAlertButton } from "@/features/developer/developer-actions";
import { IncidentList, LatencyCard, ObservationCard, RegionalList } from "@/features/developer/metric-cards";
import { DeveloperUnavailable } from "@/features/developer/unavailable";
import { requireAccount } from "@/lib/account/guard";
import { getDevAlerts, getDevDashboard } from "@/lib/api/developer";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Developer dashboard",
  robots: { index: false, follow: false },
};

export default async function DeveloperDashboardPage() {
  const { unavailable } = await requireAccount("/developers/dashboard");
  if (unavailable) {
    return (
      <DeveloperUnavailable
        title="Developer dashboard"
        description="Availability and latency stay unavailable until the API is connected."
      />
    );
  }

  const token = await readSessionToken();
  const loaded = token ? await getDevDashboard(token) : null;
  const dash = loaded && loaded.ok ? loaded.dashboard : null;
  const alerts = token ? await getDevAlerts(token) : null;
  const rules = alerts && alerts.ok ? alerts.rules : [];

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Developers"
        title="Workspace dashboard"
        description="These figures use stored worker checks only. Empty series stay unmeasured."
      />
      <PageContainer className="space-y-6 py-10">
        {!dash ? (
          <DevelopmentBanner description={loaded && !loaded.ok ? loaded.message : "Sign in required."} />
        ) : (
          <>
            <p className="text-sm text-muted-foreground">{dash.summary}</p>
            <div className="grid gap-4 md:grid-cols-2">
              <ObservationCard title="Availability" item={dash.availability} />
              <LatencyCard item={dash.latency} />
              <Card>
                <CardHeader>
                  <CardTitle>Incidents</CardTitle>
                </CardHeader>
                <CardContent>
                  <IncidentList items={dash.incidents} />
                </CardContent>
              </Card>
              <Card>
                <CardHeader>
                  <CardTitle>Regional performance</CardTitle>
                </CardHeader>
                <CardContent>
                  <RegionalList items={dash.regionalPerformance} />
                </CardContent>
              </Card>
            </div>
            <Card>
              <CardHeader>
                <CardTitle>Alert rules</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {rules.length ? (
                  <ul className="space-y-2 text-sm">
                    {rules.map((rule) => (
                      <li key={rule.id} className="flex flex-wrap items-start justify-between gap-2">
                        <div>
                          <p className="font-medium">
                            {rule.kind} · threshold {rule.threshold}
                          </p>
                          <p className="text-muted-foreground">
                            Delivered: {rule.deliveredCount}. {rule.summary}
                          </p>
                        </div>
                        <DeleteAlertButton id={rule.id} />
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p className="text-sm text-muted-foreground">No alert rules are stored.</p>
                )}
                <CreateAlertForm />
              </CardContent>
            </Card>
          </>
        )}
      </PageContainer>
    </main>
  );
}
