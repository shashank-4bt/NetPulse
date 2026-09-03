import type { Metadata } from "next";
import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { PageHero } from "@/components/public/page-hero";
import { Button } from "@/components/ui/button";

export const metadata: Metadata = {
  title: "Privacy",
  description:
    "NetPulse collects only what diagnosis requires. No full browsing history.",
  alternates: { canonical: "/privacy" },
};

export default function PrivacyPage() {
  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Privacy"
        title="Data minimization"
        description="Signed-in accounts store email, a password hash, session metadata, and diagnoses you start. Browsing history is not collected."
      />
      <PageContainer className="max-w-2xl space-y-8 py-10 text-sm text-muted-foreground md:text-base">
        <section aria-labelledby="collect-heading" className="space-y-2">
          <h2 id="collect-heading" className="text-lg font-semibold text-foreground">
            What diagnosis needs
          </h2>
          <p>
            A public hostname or URL you submit, coarse network context, and
            probe observations. That is enough to isolate a layer. Full browsing
            history is not required and will not be requested.
          </p>
        </section>
        <section aria-labelledby="geo-heading" className="space-y-2">
          <h2 id="geo-heading" className="text-lg font-semibold text-foreground">
            Geography
          </h2>
          <p>
            Comparisons use coarse geography (region or metro), not a precise
            home address. Map views stay unavailable until aggregates exist.
          </p>
        </section>
        <section aria-labelledby="tel-heading" className="space-y-2">
          <h2 id="tel-heading" className="text-lg font-semibold text-foreground">
            Telemetry
          </h2>
          <p>
            Product telemetry is off by default and can be changed in
            account privacy settings. This website does not send diagnostic
            telemetry from the browser, and it does not record browsing history.
          </p>
        </section>
        <section aria-labelledby="retain-heading" className="space-y-2">
          <h2 id="retain-heading" className="text-lg font-semibold text-foreground">
            Retention
          </h2>
          <p>
            Raw observations should be short-lived. Aggregates may last longer.
            Exact TTLs will be published with the API, not invented here.
            Signed-in users can delete an account from account privacy settings
            after confirmation. That removes the account in this process; it
            does not claim a global erase of already-stored aggregates.
          </p>
        </section>
        <Button
          variant="outline"
          nativeButton={false}
          render={<Link href="/trust" />}
        >
          Trust
        </Button>
      </PageContainer>
    </main>
  );
}
