import type { Metadata } from "next";
import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { PageHero } from "@/components/public/page-hero";
import { Button } from "@/components/ui/button";

export const metadata: Metadata = {
  title: "About",
  description:
    "NetPulse is an Internet Health Intelligence Platform that isolates where a failure actually is.",
  alternates: { canonical: "/about" },
};

export default function AboutPage() {
  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="About"
        title="Internet Health Intelligence"
        description="The core question is simple: is it me, my network, my ISP, or the service? NetPulse exists to answer that with evidence."
      />
      <PageContainer className="space-y-8 py-10">
        <section aria-labelledby="position-heading" className="max-w-2xl space-y-3">
          <h2 id="position-heading" className="text-lg font-semibold">
            What we build
          </h2>
          <p className="text-sm text-muted-foreground md:text-base">
            NetPulse measures the path from a client toward a public service:
            device, Wi-Fi, DNS, connectivity, ISP, routing, CDN, TLS, HTTP, and
            the service itself. Each layer can be healthy, degraded, failed,
            not measured, or insufficient evidence. Skipped layers are never
            marked healthy.
          </p>
        </section>
        <section aria-labelledby="who-heading" className="max-w-2xl space-y-3">
          <h2 id="who-heading" className="text-lg font-semibold">
            Who it is for
          </h2>
          <ul className="list-disc space-y-2 pl-5 text-sm text-muted-foreground md:text-base">
            <li>People who need to know if an outage is local or remote.</li>
            <li>Developers who need evidence attached to a failing check.</li>
            <li>Operators who need organization membership, roles, and stored fleet records without invented health scores.</li>
            <li>
              Operators who need regional and ASN comparison without decorative
              dashboards.
            </li>
          </ul>
        </section>
        <div className="flex flex-wrap gap-3">
          <Button nativeButton={false} render={<Link href="/how-it-works" />}>
            How it works
          </Button>
          <Button
            variant="outline"
            nativeButton={false}
            render={<Link href="/trust" />}
          >
            Trust
          </Button>
        </div>
      </PageContainer>
    </main>
  );
}
