import type { Metadata } from "next";

import { ChartContainer } from "@/components/data/chart-container";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { getPublicHealthSnapshot } from "@/lib/api/public-health";

export const metadata: Metadata = {
  title: "Map",
  description:
    "Regional internet comparison. The map is unavailable until aggregates exist.",
  alternates: { canonical: "/map" },
};

export default function MapPage() {
  const health = getPublicHealthSnapshot();

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Map"
        title="Regional and network comparison"
        description="A map is shown only when coarse geography and measured aggregates are available. None are connected."
      />
      <PageContainer className="space-y-6 py-10">
        <DevelopmentBanner description={health.reason} />
        <ChartContainer
          title="Regional reachability"
          description="MapLibre is reserved for real aggregates. This container will not invent a heatmap."
          state="unavailable"
        />
      </PageContainer>
    </main>
  );
}
