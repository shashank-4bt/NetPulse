import type { Metadata } from "next";
import Link from "next/link";

import { EvidenceItem } from "@/components/data/evidence-item";
import { ConfidenceBadge } from "@/components/status/confidence-badge";
import { PageContainer } from "@/components/layout/page-container";
import { PageHero } from "@/components/public/page-hero";
import { Button } from "@/components/ui/button";
import { CONFIDENCE_LEVELS } from "@/domain/display";
import { CONFIDENCE_LABEL } from "@/lib/design/taxonomy";

export const metadata: Metadata = {
  title: "Trust",
  description:
    "How NetPulse separates measured facts from hypotheses, states confidence, and limits data collection.",
  alternates: { canonical: "/trust" },
};

export default function TrustPage() {
  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Trust"
        title="Evidence before explanation"
        description="NetPulse is useful only if you can tell a measurement from a guess. That distinction is part of the product, not a disclaimer in small type."
      />
      <PageContainer className="space-y-10 py-10">
        <section aria-labelledby="facts-heading" className="space-y-4">
          <h2 id="facts-heading" className="text-lg font-semibold">
            Measured facts
          </h2>
          <p className="max-w-2xl text-sm text-muted-foreground">
            A fact cites a probe, a timestamp, and an observation. If the probe
            did not run, there is no fact to display.
          </p>
          <EvidenceItem
            evidenceClass="measured_fact"
            title="Example of a fact label"
            body="This component is for display. It is not a live observation from your network."
          />
        </section>
        <section aria-labelledby="hyp-heading" className="space-y-4">
          <h2 id="hyp-heading" className="text-lg font-semibold">
            Inferred hypotheses
          </h2>
          <p className="max-w-2xl text-sm text-muted-foreground">
            Observed evidence suggests a cause. Hypotheses use dashed treatment
            and an explicit label so they cannot be mistaken for measurements.
          </p>
          <EvidenceItem
            evidenceClass="inferred_hypothesis"
            title="Example of a hypothesis label"
            body="Observed evidence suggests a possible ISP or routing issue. That sentence is a hypothesis, not a measurement."
          />
        </section>
        <section aria-labelledby="conf-heading" className="space-y-3">
          <h2 id="conf-heading" className="text-lg font-semibold">
            Confidence
          </h2>
          <p className="max-w-2xl text-sm text-muted-foreground">
            Confidence is shown as text and a mark, not color alone. Incomplete
            evidence stays low or is omitted.
          </p>
          <div className="flex flex-wrap gap-2">
            {CONFIDENCE_LEVELS.map((level) => (
              <ConfidenceBadge key={level} confidence={level} />
            ))}
          </div>
          <p className="text-sm text-muted-foreground">
            Labels: {CONFIDENCE_LEVELS.map((level) => CONFIDENCE_LABEL[level]).join(", ")}.
          </p>
        </section>
        <section aria-labelledby="limits-heading" className="max-w-2xl space-y-3">
          <h2 id="limits-heading" className="text-lg font-semibold">
            Limitations
          </h2>
          <ul className="list-disc space-y-2 pl-5 text-sm text-muted-foreground">
            <li>Workers and stores are not connected in this deployment.</li>
            <li>No page on this site is a live outage or ISP score.</li>
            <li>A successful check of one hostname is not a global health claim.</li>
            <li>NetPulse will not accuse malware or execute destructive fixes.</li>
          </ul>
        </section>
        <section aria-labelledby="privacy-heading" className="max-w-2xl space-y-3">
          <h2 id="privacy-heading" className="text-lg font-semibold">
            Privacy, telemetry, minimization
          </h2>
          <p className="text-sm text-muted-foreground">
            Diagnosis needs a target and coarse network context. It does not
            need your browsing history. Telemetry, when added, will be explicit
            and off by default unless you opt in. Retention of raw observations
            will stay short.
          </p>
          <div className="flex flex-wrap gap-3">
            <Button
              variant="outline"
              nativeButton={false}
              render={<Link href="/privacy" />}
            >
              Privacy
            </Button>
            <Button
              variant="outline"
              nativeButton={false}
              render={<Link href="/security" />}
            >
              Security
            </Button>
          </div>
        </section>
      </PageContainer>
    </main>
  );
}
