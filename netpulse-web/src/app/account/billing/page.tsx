import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { requireAccount } from "@/lib/account/guard";
import { getBilling } from "@/lib/api/account";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Billing",
  robots: { index: false, follow: false },
};

export default async function BillingPage() {
  const { unavailable } = await requireAccount("/account/billing");
  if (unavailable) {
    return (
      <main id="main-content" className="flex-1">
        <PageHero eyebrow="Account" title="Billing" description="Billing needs the API." />
        <PageContainer className="py-10">
          <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set. No sample invoices are shown." />
        </PageContainer>
      </main>
    );
  }
  const token = await readSessionToken();
  const loaded = token ? await getBilling(token) : null;
  const billing = loaded && loaded.ok ? loaded.billing : null;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Account"
        title="Billing"
        description="Organizations and invoices are not enabled. This page does not invent a plan or payment history."
      />
      <PageContainer className="max-w-xl space-y-4 py-10 text-sm text-muted-foreground">
        {billing ? (
          <>
            <p>{billing.summary}</p>
            <p>Has billing account: {billing.hasAccount ? "yes" : "no"}.</p>
            <p>Invoices stored: {billing.invoices.length}.</p>
          </>
        ) : (
          <DevelopmentBanner description={loaded && !loaded.ok ? loaded.message : "Billing unavailable."} />
        )}
      </PageContainer>
    </main>
  );
}
