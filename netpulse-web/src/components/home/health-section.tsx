import Link from "next/link";

import { MetricCard } from "@/components/data/metric-card";
import { PageContainer } from "@/components/layout/page-container";
import { SectionContainer } from "@/components/layout/section-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { SectionHeading } from "@/components/public/section-heading";
import { Button } from "@/components/ui/button";
import { getPublicHealthSnapshot } from "@/lib/api/public-health";

export async function HealthSection() {
  const health = await getPublicHealthSnapshot();

  return (
    <SectionContainer labelledBy="health-heading" className="bg-muted/20">
      <PageContainer>
        <SectionHeading
          id="health-heading"
          eyebrow="Global internet health"
          title="Regional comparison requires measurements"
          description={health.reason}
        />
        <div className="mt-6">
          <DevelopmentBanner />
        </div>
        <div className="mt-4 grid gap-4 md:grid-cols-3">
          <MetricCard
            title="Reachability"
            description="Share of measured paths that completed"
            value={null}
            caption="Not measured. Backend unavailable."
          />
          <MetricCard
            title="DNS"
            description="Resolver success across observed networks"
            value={null}
            caption="Not measured. Backend unavailable."
          />
          <MetricCard
            title="TLS"
            description="Handshake completion on public targets"
            value={null}
            caption="Not measured. Backend unavailable."
          />
        </div>
        <Button
          className="mt-6"
          variant="outline"
          nativeButton={false}
          render={<Link href="/map" />}
        >
          Open map
        </Button>
      </PageContainer>
    </SectionContainer>
  );
}
