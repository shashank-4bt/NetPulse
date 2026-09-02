import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { AlertsForm } from "@/features/account/account-actions";
import { requireAccount } from "@/lib/account/guard";
import { getAlerts } from "@/lib/api/account";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Alerts",
  robots: { index: false, follow: false },
};

export default async function AlertsPage() {
  const { unavailable } = await requireAccount("/account/alerts");
  if (unavailable) {
    return (
      <main id="main-content" className="flex-1">
        <PageHero eyebrow="Account" title="Alerts" description="Alert preferences need the API." />
        <PageContainer className="py-10">
          <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set." />
        </PageContainer>
      </main>
    );
  }
  const token = await readSessionToken();
  const loaded = token ? await getAlerts(token) : null;
  const alerts = loaded && loaded.ok ? loaded.alerts : null;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Account"
        title="Alerts"
        description="Preferences only. Notification delivery is not configured, so delivered count stays at the stored value."
      />
      <PageContainer className="max-w-xl space-y-4 py-10">
        {alerts ? (
          <>
            <p className="text-sm text-muted-foreground">
              {alerts.summary} Delivered: {alerts.deliveredCount}.
            </p>
            <AlertsForm emailEnabled={alerts.emailEnabled} incidentAlerts={alerts.incidentAlerts} />
          </>
        ) : (
          <DevelopmentBanner description={loaded && !loaded.ok ? loaded.message : "Alerts unavailable."} />
        )}
      </PageContainer>
    </main>
  );
}
