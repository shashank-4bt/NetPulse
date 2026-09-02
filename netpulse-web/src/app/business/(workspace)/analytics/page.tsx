import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { BusinessScreen, BusinessUnavailable, NoOrganization } from "@/features/business/unavailable";
import { ObservationCard, LatencyCard } from "@/features/developer/metric-cards";
import { ANALYTICS_FILTERS } from "@/domain/business";
import { getOrgAnalytics } from "@/lib/api/business";
import { requireBusiness } from "@/lib/business/page";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Organization analytics",
  robots: { index: false, follow: false },
};

type PageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

function first(value: string | string[] | undefined): string {
  return Array.isArray(value) ? (value[0] ?? "") : (value ?? "");
}

export default async function BusinessAnalyticsPage({ searchParams }: PageProps) {
  const ctx = await requireBusiness("/business/analytics");
  if (ctx.unavailable) {
    return <BusinessUnavailable title="Analytics" description="Organization analytics need the API." />;
  }
  if (!ctx.organization || !ctx.token) {
    return <NoOrganization />;
  }
  const params = await searchParams;
  const filters: Record<string, string> = {};
  for (const key of ANALYTICS_FILTERS) {
    const value = first(params[key]).trim();
    if (value) {
      filters[key] = value;
    }
  }
  const loaded = await getOrgAnalytics(ctx.token, ctx.organization.id, filters);

  return (
    <BusinessScreen
      title="Analytics"
      description="Filters apply to stored organization checks only. Empty matches stay unmeasured."
    >
      <Card>
        <CardHeader>
          <CardTitle>Filters</CardTitle>
        </CardHeader>
        <CardContent>
          <form method="get" className="grid gap-3 md:grid-cols-2">
            {ANALYTICS_FILTERS.map((key) => (
              <div key={key}>
                <label htmlFor={key} className="text-sm font-medium">
                  {key}
                </label>
                <Input id={key} name={key} defaultValue={filters[key] ?? ""} className="mt-1" />
              </div>
            ))}
            <div className="md:col-span-2">
              <Button type="submit">Apply filters</Button>
            </div>
          </form>
        </CardContent>
      </Card>
      {!loaded.ok ? (
        <DevelopmentBanner description={loaded.message} />
      ) : (
        <>
          <p className="text-sm text-muted-foreground">{loaded.analytics.summary}</p>
          <p className="text-sm text-muted-foreground">Matching samples: {loaded.analytics.sampleCount}</p>
          <div className="grid gap-4 md:grid-cols-2">
            <ObservationCard title="Availability" item={loaded.analytics.availability} />
            <LatencyCard item={loaded.analytics.latency} />
            <Card>
              <CardHeader>
                <CardTitle>Matching incidents</CardTitle>
              </CardHeader>
              <CardContent>
                {loaded.analytics.incidents.length ? (
                  <ul className="space-y-2 text-sm">
                    {loaded.analytics.incidents.map((item) => (
                      <li key={item.id}>
                        {item.title} · {item.status}
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p className="text-sm text-muted-foreground">No stored incidents match these filters.</p>
                )}
              </CardContent>
            </Card>
          </div>
        </>
      )}
    </BusinessScreen>
  );
}
