import type { Metadata } from "next";
import Link from "next/link";

import { EmptyState } from "@/components/feedback/empty-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { SaveServiceForm } from "@/features/account/account-actions";
import { requireAccount } from "@/lib/account/guard";
import { getDashboard } from "@/lib/api/account";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Dashboard",
  robots: { index: false, follow: false },
};

export default async function DashboardPage() {
  const { unavailable } = await requireAccount("/dashboard");
  if (unavailable) {
    return (
      <main id="main-content" className="flex-1">
        <PageHero
          eyebrow="Dashboard"
          title="Account dashboard"
          description="Internet Health, diagnoses, and alerts stay unavailable until the API is connected."
        />
        <PageContainer className="py-10">
          <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set. No demo account or invented health score is shown." />
        </PageContainer>
      </main>
    );
  }

  const token = await readSessionToken();
  const loaded = token ? await getDashboard(token) : null;
  const dash = loaded && loaded.ok ? loaded.dashboard : null;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Dashboard"
        title="Your workspace"
        description="This page shows your stored work and public observatory rows. It is not a live score for your home path."
      />
      <PageContainer className="space-y-6 py-10">
        {!dash ? (
          <EmptyState
            title="Dashboard unavailable"
            description={loaded && !loaded.ok ? loaded.message : "Sign in required."}
          />
        ) : (
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Internet Health</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                {dash.internetHealth}
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Network information</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                {dash.networkInfo}
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Recent diagnoses</CardTitle>
              </CardHeader>
              <CardContent>
                <ItemList
                  items={dash.diagnoses.map((item) => ({
                    href: `/reports/${item.id}`,
                    title: item.target,
                    detail: `${item.status} · ${item.createdAt}`,
                  }))}
                  empty="No diagnoses are stored for this account."
                />
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Saved services</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <ItemList
                  items={dash.savedServices.map((item) => ({
                    href: `/service/${item.slug}`,
                    title: item.slug,
                    detail: item.createdAt,
                  }))}
                  empty="No saved services."
                />
                <SaveServiceForm />
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Active incidents</CardTitle>
              </CardHeader>
              <CardContent>
                <ItemList
                  items={dash.incidents.map((item) => ({
                    href: `/incident/${item.id}`,
                    title: item.title,
                    detail: `${item.status} · ${item.startedAt}`,
                  }))}
                  empty="No stored incidents are open."
                />
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Reports</CardTitle>
              </CardHeader>
              <CardContent>
                <ItemList
                  items={dash.reports.map((item) => ({
                    href: `/dashboard/reports`,
                    title: item.target,
                    detail: item.status,
                  }))}
                  empty="No reports are stored for this account."
                />
              </CardContent>
            </Card>
            <Card className="md:col-span-2">
              <CardHeader>
                <CardTitle>Alerts</CardTitle>
              </CardHeader>
              <CardContent className="text-sm text-muted-foreground">
                {dash.alerts.summary} Delivered count: {dash.alerts.deliveredCount}.
              </CardContent>
            </Card>
          </div>
        )}
      </PageContainer>
    </main>
  );
}

function ItemList({
  items,
  empty,
}: {
  items: Array<{ href: string; title: string; detail: string }>;
  empty: string;
}) {
  if (!items.length) {
    return <p className="text-sm text-muted-foreground">{empty}</p>;
  }
  return (
    <ul className="space-y-2 text-sm">
      {items.map((item) => (
        <li key={item.href + item.title}>
          <Link href={item.href} className="font-medium hover:underline">
            {item.title}
          </Link>
          <p className="text-muted-foreground">{item.detail}</p>
        </li>
      ))}
    </ul>
  );
}
