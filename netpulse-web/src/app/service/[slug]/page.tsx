import type { Metadata } from "next";
import { notFound } from "next/navigation";

import { ErrorState } from "@/components/feedback/error-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { ServiceIntelligenceView } from "@/features/observatory/service-intelligence-view";
import { getServiceBySlug, getServiceSlugs } from "@/lib/content/services";
import { loadServicePage } from "@/lib/observatory/load";

export const dynamic = "force-dynamic";

type ServicePageProps = {
  params: Promise<{ slug: string }>;
};

export function generateStaticParams() {
  return getServiceSlugs().map((slug) => ({ slug }));
}

export async function generateMetadata({
  params,
}: ServicePageProps): Promise<Metadata> {
  const { slug } = await params;
  const loaded = await loadServicePage(slug);
  const service = loaded?.catalog ?? getServiceBySlug(slug);
  if (!service) {
    return { title: "Service not found", robots: { index: false, follow: false } };
  }
  const measured = loaded?.intelligence.availability.measured;
  return {
    title: service.name,
    description: measured
      ? `${service.name} service intelligence from stored NetPulse measurements. Worker vantage is not a user-path diagnosis.`
      : `${service.name} diagnosis path. Live health, availability, and latency are not measured.`,
    alternates: { canonical: `/service/${service.slug}` },
  };
}

export default async function ServiceDetailPage({ params }: ServicePageProps) {
  const { slug } = await params;
  const loaded = await loadServicePage(slug);
  if (!loaded) {
    notFound();
  }

  const chartState =
    loaded.state === "unavailable"
      ? "unavailable"
      : loaded.intelligence.availability.sampleCount === 0
        ? "empty"
        : "insufficient_evidence";

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow={loaded.catalog.category}
        title={loaded.catalog.name}
        description={`${loaded.catalog.summary} Current state, health, and charts stay empty until measured series exist.`}
      />
      <PageContainer className="space-y-6 py-10">
        <DevelopmentBanner title="Service intelligence" description={loaded.reason} />
        {loaded.state === "error" ? (
          <ErrorState
            title="Service intelligence unavailable"
            description={loaded.reason}
          />
        ) : null}
        <ServiceIntelligenceView
          catalog={loaded.catalog}
          intelligence={loaded.intelligence}
          incidents={loaded.incidents}
          chartState={chartState}
        />
      </PageContainer>
    </main>
  );
}
