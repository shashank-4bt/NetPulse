import type { Metadata } from "next";
import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Button } from "@/components/ui/button";
import { getCurrentUser } from "@/lib/auth/session";
import { isApiConfigured } from "@/lib/api/backend";

export const metadata: Metadata = {
  title: "Business",
  description:
    "Signed-in organizations can store members, devices, networks, and reports. Health stays unmeasured until checks are stored.",
  alternates: { canonical: "/business" },
};

export default async function BusinessPage() {
  const user = await getCurrentUser();

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Business"
        title="Organization internet health"
        description="A signed-in membership can store devices, networks, services, and incidents for one organization. Empty series stay unmeasured. Cross-organization reads are rejected."
      />
      <PageContainer className="space-y-8 py-10">
        {!isApiConfigured() ? (
          <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set. The business workspace is unavailable. No sample members or invented availability percentages are shown." />
        ) : user ? (
          <div className="flex flex-wrap gap-3">
            <Button nativeButton={false} render={<Link href="/business/dashboard" />}>
              Open organization
            </Button>
          </div>
        ) : (
          <div className="flex flex-wrap gap-3">
            <Button nativeButton={false} render={<Link href="/login?next=/business/dashboard" />}>
              Sign in to an organization
            </Button>
          </div>
        )}
        <section aria-labelledby="contract-heading" className="max-w-2xl space-y-3">
          <h2 id="contract-heading" className="text-lg font-semibold">
            Organization contract
          </h2>
          <ul className="list-disc space-y-2 pl-5 text-sm text-muted-foreground">
            <li>Versioned HTTP JSON under /v1/orgs, plus a cookie BFF at /api/orgs.</li>
            <li>Roles are Owner, Admin, Security Admin, Developer, Analyst, Viewer, and Billing Admin.</li>
            <li>Foreign organization ids return 404. Missing permission on a membership returns 403.</li>
            <li>Organization API keys use the npo_ prefix, are hashed, and cannot mint other keys.</li>
            <li>Availability and percentiles stay unset until stored checks exist. This page does not claim live fleet health.</li>
          </ul>
        </section>
        <Button variant="outline" nativeButton={false} render={<Link href="/how-it-works" />}>
          Method
        </Button>
      </PageContainer>
    </main>
  );
}
