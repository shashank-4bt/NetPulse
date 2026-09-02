import type { Metadata } from "next";
import { notFound } from "next/navigation";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DeleteMonitorButton, RunMonitorButton, UpdateMonitorForm } from "@/features/developer/developer-actions";
import { DeveloperUnavailable } from "@/features/developer/unavailable";
import { requireAccount } from "@/lib/account/guard";
import { getDevMonitor } from "@/lib/api/developer";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Monitor",
  robots: { index: false, follow: false },
};

export default async function DeveloperMonitorPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const { unavailable } = await requireAccount(`/developers/monitors/${id}`);
  if (unavailable) {
    return <DeveloperUnavailable title="Monitor" description="Monitor detail needs the API." />;
  }
  const token = await readSessionToken();
  const loaded = token ? await getDevMonitor(token, id) : null;
  if (loaded && !loaded.ok && loaded.status === 404) {
    notFound();
  }
  const monitor = loaded && loaded.ok ? loaded.monitor : null;
  const checks = loaded && loaded.ok ? loaded.checks : [];

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Developers"
        title={monitor?.name ?? "Monitor"}
        description="Stored configuration and worker vantage checks. Status is unmeasured until a check is stored."
      />
      <PageContainer className="space-y-6 py-10">
        {!monitor ? (
          <DevelopmentBanner description={loaded && !loaded.ok ? loaded.message : "Monitor unavailable."} />
        ) : (
          <>
            <p className="text-sm text-muted-foreground">
              {monitor.type} · {monitor.target} · {monitor.status} · frequency {monitor.frequencySeconds}s ·
              timeout {monitor.timeoutSeconds}s · regions {monitor.regions.join(", ") || "none requested"}
            </p>
            <p className="text-sm text-muted-foreground">{monitor.summary}</p>
            <div className="flex flex-wrap gap-3">
              <RunMonitorButton id={monitor.id} />
              <DeleteMonitorButton id={monitor.id} />
            </div>
            <Card>
              <CardHeader>
                <CardTitle>Stored checks</CardTitle>
              </CardHeader>
              <CardContent>
                {checks.length === 0 ? (
                  <p className="text-sm text-muted-foreground">No checks are stored for this monitor.</p>
                ) : (
                  <ul className="space-y-2 text-sm">
                    {checks.map((check) => (
                      <li key={check.id}>
                        <p className="font-medium">
                          {check.ok ? "Succeeded" : "Failed"} · {check.region} · {check.at}
                        </p>
                        <p className="text-muted-foreground">
                          {check.latencyMs === null ? "No latency sample" : `${check.latencyMs} ms`} · {check.summary}
                        </p>
                      </li>
                    ))}
                  </ul>
                )}
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Edit</CardTitle>
              </CardHeader>
              <CardContent>
                <UpdateMonitorForm
                  id={monitor.id}
                  name={monitor.name}
                  target={monitor.target}
                  type={monitor.type}
                  regions={monitor.regions}
                  frequencySeconds={monitor.frequencySeconds}
                  timeoutSeconds={monitor.timeoutSeconds}
                />
              </CardContent>
            </Card>
          </>
        )}
      </PageContainer>
    </main>
  );
}
