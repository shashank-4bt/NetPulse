import type { Metadata } from "next";
import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Button } from "@/components/ui/button";

export const metadata: Metadata = {
  title: "Developers",
  description:
    "Monitor internet path health with evidence-backed checks. The public API is not connected yet.",
  alternates: { canonical: "/developers" },
};

export default function DevelopersPage() {
  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Developers"
        title="Attach evidence to a failing check"
        description="When the API ships, you will create diagnostic runs, read layer outcomes, and distinguish facts from hypotheses."
      />
      <PageContainer className="space-y-8 py-10">
        <DevelopmentBanner description="There is no public API, SDK, or webhook endpoint in this release. No sample keys are shown." />
        <section aria-labelledby="contract-heading" className="max-w-2xl space-y-3">
          <h2 id="contract-heading" className="text-lg font-semibold">
            Intended contract
          </h2>
          <ul className="list-disc space-y-2 pl-5 text-sm text-muted-foreground">
            <li>Versioned HTTP JSON under /v1.</li>
            <li>Run status includes queued, running, complete, partial, failed, unavailable, and insufficient_evidence.</li>
            <li>Every evidence item is classified as a measured fact, inferred hypothesis, or recommendation.</li>
            <li>Targets are validated against SSRF rules before any worker connects.</li>
          </ul>
        </section>
        <Button
          variant="outline"
          nativeButton={false}
          render={<Link href="/how-it-works" />}
        >
          Method
        </Button>
      </PageContainer>
    </main>
  );
}
