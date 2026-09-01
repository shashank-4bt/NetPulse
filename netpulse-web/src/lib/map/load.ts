import type { MapAggregates, MapQuery } from "@/domain/map";
import { emptyMapAggregates } from "@/domain/map";
import { parseMapQuery } from "@/features/map/query";
import { stripPreciseCoordinates } from "@/features/map/privacy";
import { payloadHasRawMeasurements, sanitizeAggregates } from "@/features/map/performance";
import { getBackendMapAggregates, isApiConfigured } from "@/lib/api/backend";

export type MapLoadState = "unavailable" | "empty" | "ready" | "error";

export type LoadedMap = {
  query: MapQuery;
  aggregates: MapAggregates;
  state: MapLoadState;
  reason: string;
};

export async function loadMapAggregates(
  searchParams: Record<string, string | string[] | undefined>
): Promise<LoadedMap> {
  const query = parseMapQuery(searchParams);
  if (!isApiConfigured()) {
    return {
      query,
      aggregates: emptyMapAggregates(
        "NETPULSE_API_BASE_URL is not set. The map will not invent country health or locations."
      ),
      state: "unavailable",
      reason:
        "NETPULSE_API_BASE_URL is not set. The map will not invent country health or locations.",
    };
  }
  const result = await getBackendMapAggregates({
    level: query.level,
    parent: query.parent,
    service: query.service,
    q: query.q,
    layers: query.layers,
  });
  if (!result.ok) {
    return {
      query,
      aggregates: emptyMapAggregates(result.message),
      state: "error",
      reason: result.message,
    };
  }
  if (payloadHasRawMeasurements(result.aggregates)) {
    return {
      query,
      aggregates: emptyMapAggregates(
        "The map response included raw measurements and was rejected."
      ),
      state: "error",
      reason: "The map response included raw measurements and was rejected.",
    };
  }
  const aggregates = sanitizeAggregates({
    ...result.aggregates,
    cells: result.aggregates.cells.map(stripPreciseCoordinates),
  });
  const empty =
    aggregates.cells.length === 0 && aggregates.incidentRefs.length === 0;
  return {
    query,
    aggregates,
    state: empty ? "empty" : "ready",
    reason: aggregates.reason,
  };
}
