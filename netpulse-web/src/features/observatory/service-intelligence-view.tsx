import Link from "next/link";

import { ChartContainer } from "@/components/data/chart-container";
import { MetricCard } from "@/components/data/metric-card";
import { EmptyState } from "@/components/feedback/empty-state";
import { HealthScore } from "@/components/status/health-score";
import { LayerStatusBadge } from "@/components/status/layer-status-badge";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import type { ServiceCatalogEntry } from "@/lib/content/services";
import type { LayerStatus } from "@/domain/display";
import type {
  PublicIncidentRecord,
  ServiceIntelligence,
} from "@/domain/observatory";
import { sampleCaption } from "@/features/observatory/empty-intelligence";
import { getKnownTargets } from "@/lib/content/known-targets";

type ServiceIntelligenceViewProps = {
  catalog: ServiceCatalogEntry;
  intelligence: ServiceIntelligence;
  incidents: PublicIncidentRecord[];
  chartState: "empty" | "unavailable" | "insufficient_evidence";
};

export function ServiceIntelligenceView({
  catalog,
  intelligence,
  incidents,
  chartState,
}: ServiceIntelligenceViewProps) {
  const hostname =
    getKnownTargets().find((target) => target.slug === catalog.slug)?.hostname ??
    catalog.slug;
  const layerStatus = toLayerStatus(intelligence.currentState);

  return (
    <div className="space-y-10">
      <section
        aria-labelledby="overview-heading"
        className="grid gap-4 md:grid-cols-2"
      >
        <article className="rounded-lg border border-border bg-card p-4">
          <h2 id="overview-heading" className="text-sm font-medium">
            Current state
          </h2>
          <div className="mt-3">
            <LayerStatusBadge status={layerStatus} />
          </div>
          <p className="mt-3 text-sm text-muted-foreground">
            State is shown only from stored measurements. Unknown is not treated
            as healthy.
          </p>
        </article>
        <HealthScore
          value={intelligence.health}
          confidence={null}
          label="Health"
        />
      </section>

      <section aria-labelledby="updated-heading">
        <h2 id="updated-heading" className="text-lg font-semibold">
          Last updated
        </h2>
        <p className="mt-2 font-mono text-sm text-muted-foreground">
          {intelligence.lastUpdated ?? "No measurement timestamp is stored."}
        </p>
      </section>

      <section
        aria-labelledby="metrics-heading"
        className="grid gap-4 md:grid-cols-3"
      >
        <h2 id="metrics-heading" className="sr-only">
          Availability, latency, and errors
        </h2>
        <MetricCard
          title="Availability"
          description={intelligence.availability.summary}
          value={intelligence.availability.measured ? intelligence.availability.value : null}
          unit={intelligence.availability.unit ?? undefined}
          caption={sampleCaption(intelligence.availability.sampleCount)}
        />
        <MetricCard
          title="Latency"
          description={intelligence.latency.summary}
          value={intelligence.latency.measured ? intelligence.latency.value : null}
          unit={intelligence.latency.unit ?? undefined}
          caption={sampleCaption(intelligence.latency.sampleCount)}
        />
        <MetricCard
          title="Errors"
          description={intelligence.errors.summary}
          value={intelligence.errors.measured ? intelligence.errors.value : null}
          unit={intelligence.errors.unit ?? undefined}
          caption={sampleCaption(intelligence.errors.sampleCount)}
        />
      </section>

      <section
        aria-labelledby="charts-heading"
        className="grid gap-4 lg:grid-cols-3"
      >
        <h2 id="charts-heading" className="sr-only">
          Historical performance charts
        </h2>
        <ChartContainer
          title="Availability"
          description="Plotted only from stored availability samples."
          state={chartState}
        />
        <ChartContainer
          title="Latency"
          description="Plotted only from stored latency samples."
          state={chartState}
        />
        <ChartContainer
          title="Errors"
          description="Plotted only from stored error samples."
          state={chartState}
        />
      </section>

      <SliceList
        title="Regional health"
        emptyTitle="No regional series"
        emptyDescription="Regional health appears only when coarse geography is attached to stored measurements."
        slices={intelligence.regionalHealth}
      />
      <SliceList
        title="Network health"
        emptyTitle="No network series"
        emptyDescription="Network or ASN health appears only when a network identifier is attached to stored measurements."
        slices={intelligence.networkHealth}
      />

      <section aria-labelledby="incidents-heading">
        <h2 id="incidents-heading" className="text-lg font-semibold">
          Recent incidents
        </h2>
        {incidents.length === 0 ? (
          <EmptyState
            className="mt-3"
            title="No stored incidents"
            description="An empty list is not a claim that this service is healthy."
          />
        ) : (
          <ul className="mt-3 space-y-2">
            {incidents.map((incident) => (
              <li key={incident.id}>
                <Link
                  href={`/incident/${incident.id}`}
                  className="text-sm font-medium hover:underline"
                >
                  {incident.title}
                </Link>
                <p className="text-sm text-muted-foreground">
                  Started {incident.startedAt}. Sample count {incident.sampleCount}.
                </p>
              </li>
            ))}
          </ul>
        )}
      </section>

      <section aria-labelledby="layers-heading">
        <h2 id="layers-heading" className="text-lg font-semibold">
          Layers this catalog entry covers
        </h2>
        <ul className="mt-3 flex flex-wrap gap-2">
          {catalog.layers.map((layer) => (
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
    </div>
  );
}

type SliceListProps = {
  title: string;
  emptyTitle: string;
  emptyDescription: string;
  slices: ServiceIntelligence["regionalHealth"];
};

function SliceList({
  title,
  emptyTitle,
  emptyDescription,
  slices,
}: SliceListProps) {
  const headingId = `${title.toLowerCase().replace(/\s+/g, "-")}-heading`;
  return (
    <section aria-labelledby={headingId}>
      <h2 id={headingId} className="text-lg font-semibold">
        {title}
      </h2>
      {slices.length === 0 ? (
        <EmptyState
          className="mt-3"
          title={emptyTitle}
          description={emptyDescription}
        />
      ) : (
        <ul className="mt-3 grid gap-3 sm:grid-cols-2">
          {slices.map((slice) => (
            <li
              key={slice.id}
              className="rounded-lg border border-border bg-card p-4"
            >
              <div className="flex items-center justify-between gap-2">
                <p className="text-sm font-medium">{slice.label}</p>
                <LayerStatusBadge status={slice.status} />
              </div>
              <p className="mt-2 text-sm text-muted-foreground">{slice.summary}</p>
              <p className="mt-2 font-mono text-xs text-muted-foreground">
                {sampleCaption(slice.sampleCount)}
              </p>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}

function toLayerStatus(state: ServiceIntelligence["currentState"]): LayerStatus {
  if (
    state === "not_measured" ||
    state === "insufficient_evidence" ||
    state === "degraded"
  ) {
    return state;
  }
  return "not_measured";
}
