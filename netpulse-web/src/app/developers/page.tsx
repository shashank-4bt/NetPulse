import type { Metadata } from "next";
import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Button } from "@/components/ui/button";
import { getCurrentUser } from "@/lib/auth/session";
import { isApiConfigured } from "@/lib/api/backend";

export const metadata: Metadata = {
  title: "Developers",
  description:
    "Signed-in workspaces can store HTTP, DNS, and TLS monitors. Empty series stay unmeasured. No sample API keys are shown here.",
  alternates: { canonical: "/developers" },
};

export default async function DevelopersPage() {
  const user = await getCurrentUser();

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Developers"
        title="Attach evidence to a failing check"
        description="A signed-in workspace can store monitors, hashed API keys, and signed webhooks. SLA figures stay unmeasured until checks are stored."
      />
      <PageContainer className="space-y-8 py-10">
        {!isApiConfigured() ? (
          <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set. The developer workspace is unavailable. No sample keys or invented SLA percentages are shown." />
        ) : user ? (
          <div className="flex flex-wrap gap-3">
            <Button nativeButton={false} render={<Link href="/developers/dashboard" />}>
              Open workspace
            </Button>
          </div>
        ) : (
          <div className="flex flex-wrap gap-3">
            <Button nativeButton={false} render={<Link href="/login?next=/developers/dashboard" />}>
              Sign in to the workspace
            </Button>
          </div>
        )}
        <section aria-labelledby="contract-heading" className="max-w-2xl space-y-3">
          <h2 id="contract-heading" className="text-lg font-semibold">
            Workspace contract
          </h2>
          <ul className="list-disc space-y-2 pl-5 text-sm text-muted-foreground">
            <li>Versioned HTTP JSON under /v1/dev, plus a cookie BFF at /api/dev.</li>
            <li>Monitor types are HTTP, DNS, and TLS. Regions are requested vantages.</li>
            <li>API keys are hashed. The raw secret is returned once on create or rotate.</li>
            <li>Webhook URLs must be HTTPS and pass SSRF checks. Payloads include event id, timestamp, signature, and an idempotency key.</li>
            <li>Availability and percentiles stay unset until enough stored checks exist. This page does not claim live multi-region SLA.</li>
          </ul>
        </section>
        <Button variant="outline" nativeButton={false} render={<Link href="/how-it-works" />}>
          Method
        </Button>
      </PageContainer>
    </main>
  );
}
