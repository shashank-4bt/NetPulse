import type { Metadata } from "next";
import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { EvidencePath } from "@/components/public/evidence-path";
import { IsolationPath } from "@/components/public/isolation-path";
import { PageHero } from "@/components/public/page-hero";
import { Button } from "@/components/ui/button";

export const metadata: Metadata = {
  title: "How it works",
  description:
    "NetPulse isolates failures from device to service, then explains them as evidence, hypothesis, confidence, and recommendation.",
  alternates: { canonical: "/how-it-works" },
};

export default function HowItWorksPage() {
  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Method"
        title="Measure the path. Isolate the layer."
        description="Diagnosis follows a fixed order. Unmeasured layers stay unlabeled. Inferences never replace facts."
      />
      <PageContainer className="space-y-10 py-10">
        <section aria-labelledby="layers-heading">
          <h2 id="layers-heading" className="text-lg font-semibold">
            Isolation path
          </h2>
          <p className="mt-2 max-w-2xl text-sm text-muted-foreground">
            Device → Wi-Fi → DNS → Connectivity → ISP → Routing → CDN → TLS →
            HTTP → Service.
          </p>
          <div className="mt-6">
            <IsolationPath />
          </div>
        </section>
        <section aria-labelledby="explain-heading">
          <h2 id="explain-heading" className="text-lg font-semibold">
            Explain only from evidence
          </h2>
          <p className="mt-2 max-w-2xl text-sm text-muted-foreground">
            Evidence → Hypothesis → Confidence → Recommendation. If the run is
            partial, the result is insufficient evidence — not a fabricated
            root cause.
          </p>
          <div className="mt-6">
            <EvidencePath />
          </div>
        </section>
        <Button nativeButton={false} render={<Link href="/#diagnose" />}>
          Check My Internet
        </Button>
      </PageContainer>
    </main>
  );
}
