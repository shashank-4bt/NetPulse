import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";

import { EmptyState } from "@/components/feedback/empty-state";
import { ErrorState } from "@/components/feedback/error-state";
import { UnavailableState } from "@/components/feedback/unavailable-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Button } from "@/components/ui/button";
import { IncidentDetailView } from "@/features/observatory/incident-detail-view";
import { loadIncident } from "@/lib/observatory/load";

export const dynamic = "force-dynamic";

type IncidentPageProps = {
  params: Promise<{ id: string }>;
};

export async function generateMetadata({
  params,
}: IncidentPageProps): Promise<Metadata> {
  const { id } = await params;
  const loaded = await loadIncident(id);
  if (loaded.missing || !loaded.incident) {
    return {
      title: "Incident not found",
      description: "No stored incident matches this id. Missing is not a measurement failure.",
      robots: { index: false, follow: false },
    };
  }
  return {
    title: loaded.incident.title,
    description: `${loaded.incident.title}. Stored incident record with observed sample count ${loaded.incident.sampleCount}. Not a population impact estimate.`,
    alternates: { canonical: `/incident/${loaded.incident.id}` },
    robots: { index: true, follow: true },
  };
}

export default async function IncidentPage({ params }: IncidentPageProps) {
  const { id } = await params;
  const loaded = await loadIncident(id);
  if (loaded.missing) {
    notFound();
  }

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Incident"
        title={loaded.incident?.title ?? "Incident unavailable"}
        description={
          loaded.incident
            ? "Stored incident intelligence. Language stays observational unless isolation evidence exists."
            : loaded.reason
        }
      />
      <PageContainer className="space-y-6 py-10">
        <DevelopmentBanner title="Incident document" description={loaded.reason} />
        {loaded.state === "unavailable" ? (
          <UnavailableState title="Incident store unavailable" description={loaded.reason} />
        ) : null}
        {loaded.state === "error" ? (
          <ErrorState title="Incident failed to load" description={loaded.reason} />
        ) : null}
        {loaded.incident ? (
          <IncidentDetailView incident={loaded.incident} />
        ) : loaded.state !== "error" && loaded.state !== "unavailable" ? (
          <EmptyState
            title="Incident not found"
            description="The store does not contain this incident."
          />
        ) : null}
        <Button nativeButton={false} render={<Link href="/outages" />}>
          Back to outages
        </Button>
      </PageContainer>
    </main>
  );
}
