import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CreateOrgForm, CreateOrgMonitorForm, DeleteOrgMonitorButton } from "@/features/business/actions";
import { BusinessScreen, BusinessUnavailable } from "@/features/business/unavailable";
import { ObservationCard, RegionalList } from "@/features/developer/metric-cards";
import { requireBusiness } from "@/lib/business/page";
import { getOrgDashboard, getOrgMonitors } from "@/lib/api/business";
import { hasPermission } from "@/domain/business";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Business dashboard",
  robots: { index: false, follow: false },
};

export default async function BusinessDashboardPage() {
  const ctx = await requireBusiness("/business/dashboard");
  if (ctx.unavailable) {
    return (
      <BusinessUnavailable
        title="Organization dashboard"
        description="Overall health stays unavailable until the API is connected."
      />
    );
  }

  if (!ctx.organization || !ctx.token) {
    return (
      <BusinessScreen
        title="Create an organization"
        description="Membership is required before any organization resource can be read."
      >
        <Card>
          <CardHeader>
            <CardTitle>New organization</CardTitle>
          </CardHeader>
          <CardContent>
            <CreateOrgForm />
          </CardContent>
        </Card>
      </BusinessScreen>
    );
  }

  const org = ctx.organization;
  const loaded = await getOrgDashboard(ctx.token, org.id);
  const dash = loaded.ok ? loaded.dashboard : null;
  const monitors = await getOrgMonitors(ctx.token, org.id);
  const canCreate = hasPermission(ctx.permissions, "monitor.create");

  return (
    <BusinessScreen
      title={org.name}
      description="These figures use stored organization checks only. Empty series stay unmeasured."
    >
      {!dash ? (
        <DevelopmentBanner description={!loaded.ok ? loaded.message : "Dashboard unavailable."} />
      ) : (
        <>
          <p className="text-sm text-muted-foreground">{dash.summary}</p>
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Overall health</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                <p>{dash.overallHealth}</p>
              </CardContent>
            </Card>
            <ObservationCard title="Availability" item={dash.availability} />
            <Card>
              <CardHeader>
                <CardTitle>Incidents</CardTitle>
              </CardHeader>
              <CardContent>
                {dash.incidents.length ? (
                  <ul className="space-y-2 text-sm">
                    {dash.incidents.map((item) => (
                      <li key={item.id}>
                        <p className="font-medium">{item.title}</p>
                        <p className="text-muted-foreground">
                          {item.status} · samples {item.sampleCount}
                        </p>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p className="text-sm text-muted-foreground">No organization incidents are stored.</p>
                )}
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Affected devices</CardTitle>
              </CardHeader>
              <CardContent>
                {dash.affectedDevices.length ? (
                  <ul className="space-y-2 text-sm">
                    {dash.affectedDevices.map((item) => (
                      <li key={item.id}>
                        {item.name} · {item.region || "no region"}
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p className="text-sm text-muted-foreground">No devices are linked to stored incidents.</p>
                )}
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Regions</CardTitle>
              </CardHeader>
              <CardContent>
                <RegionalList items={dash.regions} />
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Networks</CardTitle>
              </CardHeader>
              <CardContent>
                {dash.networks.length ? (
                  <ul className="space-y-2 text-sm">
                    {dash.networks.map((item) => (
                      <li key={item.id}>
                        {item.name} · {item.asn || "no ASN"}
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p className="text-sm text-muted-foreground">No organization networks are stored.</p>
                )}
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Services</CardTitle>
              </CardHeader>
              <CardContent>
                {dash.services.length ? (
                  <ul className="space-y-2 text-sm">
                    {dash.services.map((item) => (
                      <li key={item.id}>{item.name}</li>
                    ))}
                  </ul>
                ) : (
                  <p className="text-sm text-muted-foreground">No organization services are stored.</p>
                )}
              </CardContent>
            </Card>
          </div>
        </>
      )}
      <Card>
        <CardHeader>
          <CardTitle>Organization monitors</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {!monitors.ok ? (
            <DevelopmentBanner description={monitors.message} />
          ) : monitors.monitors.length === 0 ? (
            <p className="text-sm text-muted-foreground">No organization monitors are stored.</p>
          ) : (
            <ul className="space-y-3 text-sm">
              {monitors.monitors.map((item) => (
                <li key={item.id} className="flex flex-wrap items-start justify-between gap-2">
                  <div>
                    <p className="font-medium">
                      {item.name} · {item.type} · {item.status}
                    </p>
                    <p className="text-muted-foreground">
                      {item.target} · {item.checkCount} checks
                    </p>
                    <p className="text-muted-foreground">{item.summary}</p>
                  </div>
                  {hasPermission(ctx.permissions, "monitor.delete") ? (
                    <DeleteOrgMonitorButton orgId={org.id} id={item.id} />
                  ) : null}
                </li>
              ))}
            </ul>
          )}
          {canCreate ? <CreateOrgMonitorForm orgId={org.id} /> : null}
        </CardContent>
      </Card>
    </BusinessScreen>
  );
}
