import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { DeleteAccountButton, PrivacyForm } from "@/features/account/account-actions";
import { requireAccount } from "@/lib/account/guard";
import { getPrivacy } from "@/lib/api/account";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Privacy",
  robots: { index: false, follow: false },
};

export default async function AccountPrivacyPage() {
  const { unavailable } = await requireAccount("/account/privacy");
  if (unavailable) {
    return (
      <main id="main-content" className="flex-1">
        <PageHero eyebrow="Account" title="Privacy" description="Account privacy settings need the API." />
        <PageContainer className="py-10">
          <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set." />
        </PageContainer>
      </main>
    );
  }
  const token = await readSessionToken();
  const loaded = token ? await getPrivacy(token) : null;
  const privacy = loaded && loaded.ok ? loaded.privacy : null;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Account"
        title="Privacy"
        description="What this account stores, why, how long, and how to delete it. Browsing history is not collected."
      />
      <PageContainer className="max-w-2xl space-y-6 py-10 text-sm text-muted-foreground">
        {privacy ? (
          <>
            <section className="space-y-2">
              <h2 className="text-lg font-semibold text-foreground">Data collected</h2>
              <p>{privacy.collected}</p>
            </section>
            <section className="space-y-2">
              <h2 className="text-lg font-semibold text-foreground">Purpose</h2>
              <p>{privacy.purpose}</p>
            </section>
            <section className="space-y-2">
              <h2 className="text-lg font-semibold text-foreground">Retention</h2>
              <p>{privacy.retention}</p>
            </section>
            <section className="space-y-2">
              <h2 className="text-lg font-semibold text-foreground">Deletion</h2>
              <p>{privacy.deletion}</p>
            </section>
            <section className="space-y-2">
              <h2 className="text-lg font-semibold text-foreground">Browsing history</h2>
              <p>{privacy.browsingHistory}</p>
            </section>
            <PrivacyForm telemetryOptIn={privacy.telemetryOptIn} />
            <DeleteAccountButton />
          </>
        ) : (
          <DevelopmentBanner description={loaded && !loaded.ok ? loaded.message : "Privacy settings unavailable."} />
        )}
      </PageContainer>
    </main>
  );
}
