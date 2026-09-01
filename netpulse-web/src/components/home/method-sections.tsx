import Link from "next/link";

import { EvidenceItem } from "@/components/data/evidence-item";
import { PageContainer } from "@/components/layout/page-container";
import { SectionContainer } from "@/components/layout/section-container";
import { EvidencePath } from "@/components/public/evidence-path";
import { IsolationPath } from "@/components/public/isolation-path";
import { SectionHeading } from "@/components/public/section-heading";
import { Button } from "@/components/ui/button";

export function HowItWorksPreview() {
  return (
    <SectionContainer labelledBy="how-heading">
      <PageContainer>
        <SectionHeading
          id="how-heading"
          eyebrow="How it works"
          title="Isolate the layer before blaming the service"
          description="A failure at DNS is not a website outage. NetPulse walks the path in order and leaves unmeasured layers unlabeled."
        />
        <div className="mt-6">
          <IsolationPath />
        </div>
        <Button
          className="mt-6"
          variant="outline"
          nativeButton={false}
          render={<Link href="/how-it-works" />}
        >
          Full method
        </Button>
      </PageContainer>
    </SectionContainer>
  );
}

export function EvidencePreview() {
  return (
    <SectionContainer labelledBy="evidence-heading" className="bg-muted/20">
      <PageContainer>
        <SectionHeading
          id="evidence-heading"
          eyebrow="Evidence-based diagnosis"
          title="Facts, hypotheses, and recommendations stay separate"
          description="If evidence is weak, NetPulse says insufficient evidence — not a confident story."
        />
        <div className="mt-6">
          <EvidencePath />
        </div>
        <div className="mt-4 grid gap-3 lg:grid-cols-3">
          <EvidenceItem
            evidenceClass="measured_fact"
            title="What a fact looks like"
            body="A probe completed or failed at a timestamp. The UI labels it as a measured fact."
          />
          <EvidenceItem
            evidenceClass="inferred_hypothesis"
            title="What a hypothesis looks like"
            body="Observed evidence suggests a cause. It is never styled as a measurement."
          />
          <EvidenceItem
            evidenceClass="recommendation"
            title="What a recommendation looks like"
            body="A safe next step. NetPulse does not execute dangerous remediation."
          />
        </div>
      </PageContainer>
    </SectionContainer>
  );
}

export function IntelligencePreview() {
  return (
    <SectionContainer labelledBy="intel-heading">
      <PageContainer>
        <SectionHeading
          id="intel-heading"
          eyebrow="Global internet intelligence"
          title="Compare regions and networks from the same evidence model"
          description="Observatory views appear only after many real measurements exist. Until then, maps and aggregates stay unavailable."
        />
        <Button
          className="mt-6"
          variant="outline"
          nativeButton={false}
          render={<Link href="/map" />}
        >
          Map and comparison
        </Button>
      </PageContainer>
    </SectionContainer>
  );
}

export function AudiencePreview() {
  return (
    <SectionContainer labelledBy="audience-heading" className="bg-muted/20">
      <PageContainer className="grid gap-8 md:grid-cols-2">
        <div>
          <SectionHeading
            id="audience-heading"
            eyebrow="Developer monitoring"
            title="Instrument checks without inventing uptime"
            description="Use NetPulse to see whether a failure is DNS, TLS, routing, or the application — then attach evidence to an incident."
          />
          <Button
            className="mt-4"
            variant="outline"
            nativeButton={false}
            render={<Link href="/developers" />}
          >
            Developers
          </Button>
        </div>
        <div>
          <h2 className="text-xl font-semibold tracking-tight md:text-2xl">
            Business reliability
          </h2>
          <p className="mt-2 text-sm text-muted-foreground md:text-base">
            Operators need to know if customers, an ISP, or a vendor is
            failing. Scores include confidence. Incomplete runs stay partial.
          </p>
          <Button
            className="mt-4"
            variant="outline"
            nativeButton={false}
            render={<Link href="/about" />}
          >
            About NetPulse
          </Button>
        </div>
      </PageContainer>
    </SectionContainer>
  );
}

export function TrustPreview() {
  return (
    <SectionContainer labelledBy="trust-heading">
      <PageContainer>
        <SectionHeading
          id="trust-heading"
          eyebrow="Trust and privacy"
          title="Collect only what diagnosis requires"
          description="Coarse geography, short retention, and explicit telemetry controls. No full browsing history."
        />
        <div className="mt-4 flex flex-wrap gap-3">
          <Button
            variant="outline"
            nativeButton={false}
            render={<Link href="/trust" />}
          >
            Trust
          </Button>
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
      </PageContainer>
    </SectionContainer>
  );
}

export function FinalCtaSection() {
  return (
    <SectionContainer labelledBy="final-cta-heading" className="bg-muted/20">
      <PageContainer>
        <h2
          id="final-cta-heading"
          className="text-2xl font-semibold tracking-tight"
        >
          Check a public hostname
        </h2>
        <p className="mt-2 max-w-xl text-muted-foreground">
          The form validates input now. Probes start when workers are connected.
        </p>
        <Button className="mt-6" nativeButton={false} render={<Link href="/diagnose" />}>
          Check My Internet
        </Button>
      </PageContainer>
    </SectionContainer>
  );
}
