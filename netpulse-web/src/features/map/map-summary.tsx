import Link from "next/link";

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { MAP_LEVEL_LABELS, type MapAggregates, type MapCell, type MapQuery } from "@/domain/map";
import { mapQueryString } from "@/features/map/query";
import { drillHref, drillLabel } from "@/features/map/selection";

type MapSummaryProps = {
  query: MapQuery;
  aggregates: MapAggregates;
  selected: MapCell | null;
  mapReady: boolean;
  mapError: string | null;
};

export function MapSummary({
  query,
  aggregates,
  selected,
  mapReady,
  mapError,
}: MapSummaryProps) {
  return (
    <div className="space-y-4">
      <nav aria-label="Map hierarchy">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink render={<Link href="/map" />}>World</BreadcrumbLink>
            </BreadcrumbItem>
            {query.country ? (
              <>
                <BreadcrumbSeparator />
                <BreadcrumbItem>
                  <BreadcrumbLink
                    render={
                      <Link
                        href={`/map${mapQueryString({ ...query, level: "region", parent: `country:${query.country}`, select: "" })}`}
                      />
                    }
                  >
                    {query.country}
                  </BreadcrumbLink>
                </BreadcrumbItem>
              </>
            ) : null}
            {query.region ? (
              <>
                <BreadcrumbSeparator />
                <BreadcrumbItem>
                  <BreadcrumbLink
                    render={
                      <Link
                        href={`/map${mapQueryString({ ...query, level: "network", parent: `region:${query.region}`, select: "" })}`}
                      />
                    }
                  >
                    {query.region}
                  </BreadcrumbLink>
                </BreadcrumbItem>
              </>
            ) : null}
            {query.network ? (
              <>
                <BreadcrumbSeparator />
                <BreadcrumbItem>
                  <BreadcrumbLink
                    render={
                      <Link
                        href={`/map${mapQueryString({ ...query, level: "service", parent: `network:${query.network}`, select: "" })}`}
                      />
                    }
                  >
                    {query.network}
                  </BreadcrumbLink>
                </BreadcrumbItem>
              </>
            ) : null}
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbPage>{MAP_LEVEL_LABELS[query.level]}</BreadcrumbPage>
            </BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </nav>

      <Card size="sm">
        <CardHeader>
          <CardTitle>Selection</CardTitle>
          <CardDescription>
            Click a country for a regional summary, a region for networks, a
            network for services, or a service for intelligence.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          {selected ? (
            <>
              <p className="font-medium text-foreground">{selected.label}</p>
              <p className="text-sm text-muted-foreground">{selected.summary}</p>
              <p className="font-mono text-xs text-muted-foreground">
                Observed samples: {selected.sampleCount}. Child rows: {selected.childCount}.
              </p>
              <Button nativeButton={false} render={<Link href={drillHref(selected, query)} />}>
                {drillLabel(selected)}
              </Button>
            </>
          ) : (
            <p className="text-sm text-muted-foreground">
              No row is selected. Use the table or a plottable map point.
            </p>
          )}
        </CardContent>
      </Card>

      <Card size="sm">
        <CardHeader>
          <CardTitle>Viewport load</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2 text-sm text-muted-foreground">
          <p>
            Basemap: {mapReady ? "loaded" : "loading"}. Overlay points appear only
            when a coarse centroid is stored.
          </p>
          <p>
            Precision: {aggregates.precision}. Plottable points:{" "}
            {aggregates.hasCoordinates ? "yes" : "none"}.
            {aggregates.truncated ? " This viewport result was capped." : ""}
          </p>
          {mapError ? <p role="status">{mapError}</p> : null}
        </CardContent>
      </Card>
    </div>
  );
}
