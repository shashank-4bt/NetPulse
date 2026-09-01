import type { Metadata } from "next";
import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { MapExplorer } from "@/features/map/map-explorer";
import { MapFilters } from "@/features/map/map-filters";
import { SERVICE_CATALOG } from "@/lib/content/services";
import { loadMapAggregates } from "@/lib/map/load";

export const dynamic = "force-dynamic";

type MapPageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export async function generateMetadata({
  searchParams,
}: MapPageProps): Promise<Metadata> {
  const loaded = await loadMapAggregates(await searchParams);
  return {
    title: "Map",
    description:
      loaded.state === "unavailable"
        ? "Global internet health map. Aggregates are unavailable until the API is connected."
        : "Coarse geographic aggregates only. The map does not expose precise individual locations.",
    alternates: { canonical: "/map" },
  };
}

export default async function MapPage({ searchParams }: MapPageProps) {
  const loaded = await loadMapAggregates(await searchParams);

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Map"
        title="Global internet health map"
        description="World → country → region → network/ASN → service. Layers are aggregated health overlays, not raw measurements. Precise individual locations are never shown."
        actions={
          <Link
            href="#map-data-table"
            className="text-sm text-muted-foreground underline-offset-4 hover:underline"
          >
            Skip to aggregate table
          </Link>
        }
      />
      <PageContainer className="space-y-6 py-10">
        <DevelopmentBanner description={loaded.reason} />
        <MapFilters query={loaded.query} services={SERVICE_CATALOG} />
        <MapExplorer
          query={loaded.query}
          initial={loaded.aggregates}
          state={loaded.state}
          reason={loaded.reason}
          liveEnabled={loaded.state !== "unavailable"}
        />
      </PageContainer>
    </main>
  );
}
