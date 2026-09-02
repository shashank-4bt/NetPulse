import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { GenerateReportForm } from "@/features/business/actions";
import { BusinessScreen, BusinessUnavailable, NoOrganization } from "@/features/business/unavailable";
import { ObservationCard, LatencyCard, RegionalList } from "@/features/developer/metric-cards";
import { getOrgReports } from "@/lib/api/business";
import { requireBusiness } from "@/lib/business/page";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Organization reports",
  robots: { index: false, follow: false },
};

export default async function BusinessReportsPage() {
  const ctx = await requireBusiness("/business/reports");
  if (ctx.unavailable) {
    return <BusinessUnavailable title="Reports" description="Organization reports need the API." />;
  }
  if (!ctx.organization || !ctx.token) {
    return <NoOrganization />;
  }
  const loaded = await getOrgReports(ctx.token, ctx.organization.id);

  return (
    <BusinessScreen title="Reports" description="Reports are generated from stored organization records. Empty fields stay empty.">
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Generate</CardTitle>
          </CardHeader>
          <CardContent>
            <GenerateReportForm orgId={ctx.organization.id} />
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Stored reports</CardTitle>
          </CardHeader>
          <CardContent>
            {!loaded.ok ? (
              <DevelopmentBanner description={loaded.message} />
            ) : loaded.reports.length === 0 ? (
              <p className="text-sm text-muted-foreground">No organization reports are stored.</p>
            ) : (
              <ul className="space-y-6">
                {loaded.reports.map((item) => (
                  <li key={item.id} className="space-y-3">
                    <div>
                      <p className="font-medium">{item.title}</p>
                      <p className="text-sm text-muted-foreground">
                        {item.kind} · samples {item.sampleCount} · {item.createdAt}
                      </p>
                      <p className="text-sm text-muted-foreground">{item.summary}</p>
                    </div>
                    <div className="grid gap-3 md:grid-cols-2">
                      <ObservationCard title="Availability" item={item.availability} />
                      <LatencyCard item={item.latency} />
                    </div>
                    <RegionalList items={item.regions} />
                    {item.findings.length ? (
                      <ul className="list-disc pl-5 text-sm text-muted-foreground">
                        {item.findings.map((finding) => (
                          <li key={finding}>{finding}</li>
                        ))}
                      </ul>
                    ) : (
                      <p className="text-sm text-muted-foreground">No diagnostic findings are stored on this report.</p>
                    )}
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
      </div>
    </BusinessScreen>
  );
}
