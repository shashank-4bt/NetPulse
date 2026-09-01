"use client";

import dynamic from "next/dynamic";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useMemo, useState } from "react";

import { EmptyState } from "@/components/feedback/empty-state";
import { ErrorState } from "@/components/feedback/error-state";
import { UnavailableState } from "@/components/feedback/unavailable-state";
import type { MapAggregates, MapQuery, MapViewport } from "@/domain/map";
import { MapSummary } from "@/features/map/map-summary";
import { MapTable } from "@/features/map/map-table";
import { mapAggregatesQueryString } from "@/features/map/query";
import { drillHref, selectedCell } from "@/features/map/selection";
import { stripPreciseCoordinates } from "@/features/map/privacy";
import { sanitizeAggregates } from "@/features/map/performance";

const MapCanvas = dynamic(
  () => import("@/features/map/map-canvas").then((mod) => mod.MapCanvas),
  {
    ssr: false,
    loading: () => (
      <div className="flex h-full items-center bg-muted/40 p-4 text-sm text-muted-foreground">
        Loading map canvas…
      </div>
    ),
  }
);

type MapExplorerProps = {
  query: MapQuery;
  initial: MapAggregates;
  state: "unavailable" | "empty" | "ready" | "error";
  reason: string;
  liveEnabled: boolean;
};

export function MapExplorer({
  query,
  initial,
  state,
  reason,
  liveEnabled,
}: MapExplorerProps) {
  const router = useRouter();
  const queryKey = mapAggregatesQueryString(query);
  const [live, setLive] = useState<{ key: string; data: MapAggregates } | null>(
    null
  );
  const [viewportError, setViewportError] = useState<string | null>(null);
  const [mapError, setMapError] = useState<string | null>(null);
  const [mapReady, setMapReady] = useState(false);
  const [viewport, setViewport] = useState<MapViewport | null>(null);
  const aggregates = live?.key === queryKey ? live.data : initial;

  const safeAggregates = useMemo(() => {
    const sanitized = sanitizeAggregates(aggregates);
    return {
      ...sanitized,
      cells: sanitized.cells.map(stripPreciseCoordinates),
    };
  }, [aggregates]);

  const selected = selectedCell(safeAggregates.cells, query.select);

  const onSelect = useCallback(
    (id: string) => {
      const cell = safeAggregates.cells.find((item) => item.id === id);
      if (!cell) {
        return;
      }
      router.push(drillHref(cell, query));
    },
    [query, router, safeAggregates.cells]
  );

  useEffect(() => {
    if (!liveEnabled || !viewport) {
      return;
    }
    const controller = new AbortController();
    const timer = window.setTimeout(async () => {
      try {
        const response = await fetch(
          `/api/map/aggregates?${mapAggregatesQueryString(query, viewport)}`,
          { signal: controller.signal, cache: "no-store" }
        );
        const body = (await response.json()) as {
          ok?: boolean;
          map?: MapAggregates;
          error?: { message?: string };
        };
        if (!response.ok || !body.ok || !body.map) {
          setViewportError(
            body.error?.message ??
              "Viewport aggregates failed. The table still shows the last snapshot."
          );
          return;
        }
        setLive({ key: queryKey, data: body.map });
        setViewportError(null);
      } catch (error) {
        if (controller.signal.aborted) {
          return;
        }
        setViewportError(
          error instanceof Error
            ? error.message
            : "Viewport aggregates failed. The table still shows the last snapshot."
        );
      }
    }, 280);
    return () => {
      controller.abort();
      window.clearTimeout(timer);
    };
  }, [liveEnabled, query, queryKey, viewport]);

  return (
    <div className="space-y-6">
      <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_18rem]">
        <div className="space-y-3">
          <p className="md:hidden text-sm text-muted-foreground">
            On small screens the table is the primary way to browse aggregates.
            The canvas uses two-finger gestures so the page can still scroll.
          </p>
          <div className="h-48 overflow-hidden rounded-lg border border-border md:h-[28rem]">
            <MapCanvas
              cells={safeAggregates.cells}
              selectedId={query.select}
              onSelect={onSelect}
              onViewport={setViewport}
              onLoadError={setMapError}
              onReady={() => setMapReady(true)}
            />
          </div>
          {!safeAggregates.hasCoordinates ? (
            <p className="text-sm text-muted-foreground">
              No plottable coarse points for this view. The basemap is cartography,
              not a health heatmap.
            </p>
          ) : null}
        </div>
        <MapSummary
          query={query}
          aggregates={safeAggregates}
          selected={selected}
          mapReady={mapReady}
          mapError={mapError}
        />
      </div>

      {state === "unavailable" ? (
        <UnavailableState title="Map store unavailable" description={reason} />
      ) : null}
      {state === "error" ? (
        <ErrorState title="Map aggregates failed" description={reason} />
      ) : null}
      {viewportError ? (
        <ErrorState title="Viewport query failed" description={viewportError} />
      ) : null}

      <div id="map-data-table">
        <h2 className="mb-3 text-sm font-medium">Aggregate table</h2>
        <MapTable aggregates={safeAggregates} query={query} />
      </div>

      {safeAggregates.cells.length === 0 &&
      safeAggregates.incidentRefs.length === 0 &&
      state !== "error" ? (
        <EmptyState
          title="No map aggregates"
          description="NetPulse will not color countries or place points until coarse geographic aggregates exist. Filters still apply to stored rows only."
        />
      ) : null}
    </div>
  );
}

