import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";

import { HealthScore } from "@/components/status/health-score";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { getKnownTargets } from "@/lib/content/known-targets";
import {
  getServiceBySlug,
  getServiceSlugs,
} from "@/lib/content/services";

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
  const service = getServiceBySlug(slug);
  if (!service) {
    return { title: "Service not found" };
  }
  return {
    title: service.name,
    description: `${service.name} diagnosis path. Live health is not available.`,
    alternates: { canonical: `/service/${service.slug}` },
  };
}

export default async function ServiceDetailPage({ params }: ServicePageProps) {
  const { slug } = await params;
  const service = getServiceBySlug(slug);
  if (!service) {
    notFound();
  }

  const hostname =
    getKnownTargets().find((target) => target.slug === service.slug)?.hostname ??
    service.slug;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow={service.category}
        title={service.name}
        description={`${service.summary} This page does not report live status.`}
      />
      <PageContainer className="space-y-6 py-10">
        <DevelopmentBanner
          title="Not measured"
          description={`${service.name} has no connected probe series. NetPulse will not display an invented uptime or incident.`}
        />
        <HealthScore value={null} confidence={null} label="Service health" />
        <section aria-labelledby="layers-heading">
          <h2 id="layers-heading" className="text-lg font-semibold">
            Layers this catalog entry covers
          </h2>
          <ul className="mt-3 flex flex-wrap gap-2">
            {service.layers.map((layer) => (
              <li key={layer}>
                <Badge variant="outline">{layer}</Badge>
              </li>
            ))}
          </ul>
        </section>
        <div className="flex flex-wrap gap-3">
          <Button
            nativeButton={false}
            render={
              <Link href={`/diagnose?target=${encodeURIComponent(hostname)}`} />
            }
          >
            Check My Internet
          </Button>
          <Button
            variant="outline"
            nativeButton={false}
            render={<Link href="/services" />}
          >
            All services
          </Button>
        </div>
      </PageContainer>
    </main>
  );
}
