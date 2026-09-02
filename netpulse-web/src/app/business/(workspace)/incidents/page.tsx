import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ResolveIncidentButton } from "@/features/business/actions";
import { BusinessScreen, BusinessUnavailable, NoOrganization } from "@/features/business/unavailable";
import { hasPermission } from "@/domain/business";
import { getOrgIncidents } from "@/lib/api/business";
import { requireBusiness } from "@/lib/business/page";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Organization incidents",
  robots: { index: false, follow: false },
};

export default async function BusinessIncidentsPage() {
  const ctx = await requireBusiness("/business/incidents");
  if (ctx.unavailable) {
    return <BusinessUnavailable title="Incidents" description="Organization incidents need the API." />;
  }
  if (!ctx.organization || !ctx.token) {
    return <NoOrganization />;
  }
  const org = ctx.organization;
  const loaded = await getOrgIncidents(ctx.token, org.id);
  const canManage = hasPermission(ctx.permissions, "incident.manage");

  return (
    <BusinessScreen title="Incidents" description="Incidents are opened from stored organization monitor checks only.">
      <Card>
        <CardHeader>
          <CardTitle>Stored incidents</CardTitle>
        </CardHeader>
        <CardContent>
          {!loaded.ok ? (
            <DevelopmentBanner description={loaded.message} />
          ) : loaded.incidents.length === 0 ? (
            <p className="text-sm text-muted-foreground">No organization incidents are stored.</p>
          ) : (
            <ul className="space-y-3 text-sm">
              {loaded.incidents.map((item) => (
                <li key={item.id} className="flex flex-wrap items-start justify-between gap-2">
                  <div>
                    <p className="font-medium">{item.title}</p>
                    <p className="text-muted-foreground">
                      {item.status} · samples {item.sampleCount} · {item.startedAt}
                    </p>
                    <p className="text-muted-foreground">{item.summary}</p>
                  </div>
                  {canManage && item.status !== "resolved" ? (
                    <ResolveIncidentButton orgId={org.id} id={item.id} />
                  ) : null}
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </BusinessScreen>
  );
}
