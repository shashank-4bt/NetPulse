import Link from "next/link";

import { EmptyState } from "@/components/feedback/empty-state";
import { PageContainer } from "@/components/layout/page-container";
import { SectionContainer } from "@/components/layout/section-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { SectionHeading } from "@/components/public/section-heading";
import { Button } from "@/components/ui/button";
import { getPublicHealthSnapshot } from "@/lib/api/public-health";

export function IncidentsSection() {
  const health = getPublicHealthSnapshot();

  return (
    <SectionContainer labelledBy="incidents-heading">
      <PageContainer>
        <SectionHeading
          id="incidents-heading"
          eyebrow="Current incidents"
          title="No incident feed is connected"
          description="An empty list is not the same as a healthy internet. NetPulse will not invent outages."
        />
        <div className="mt-6 space-y-4">
          <DevelopmentBanner description={health.reason} />
          <EmptyState
            title="No incidents to show"
            description="When the incident store is connected, confirmed events will appear here with evidence and confidence."
          >
            <Button
              variant="outline"
              nativeButton={false}
              render={<Link href="/outages" />}
            >
              Outages page
            </Button>
          </EmptyState>
        </div>
      </PageContainer>
    </SectionContainer>
  );
}
