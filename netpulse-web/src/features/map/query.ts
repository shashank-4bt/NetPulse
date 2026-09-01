import {
  DEFAULT_MAP_QUERY,
  MAP_LAYERS,
  MAP_LEVELS,
  type MapLayer,
  type MapLevel,
  type MapQuery,
  type MapViewport,
} from "@/domain/map";

export function parseMapQuery(
  searchParams: Record<string, string | string[] | undefined>
): MapQuery {
  const first = (key: string): string => {
    const value = searchParams[key];
    if (Array.isArray(value)) {
      return value[0]?.trim() ?? "";
    }
    return value?.trim() ?? "";
  };
  const levelRaw = first("level");
  const level = (MAP_LEVELS as readonly string[]).includes(levelRaw)
    ? (levelRaw as MapLevel)
    : DEFAULT_MAP_QUERY.level;
  return {
    level,
    parent: first("parent"),
    country: first("country"),
    region: first("region"),
    network: first("network"),
    service: first("service"),
    q: first("q"),
    layers: parseLayers(searchParams.layers),
    select: first("select"),
  };
}

export function parseLayers(value: string | string[] | undefined): MapLayer[] {
  const parts = Array.isArray(value)
    ? value.flatMap((item) => item.split(","))
    : (value ?? "").split(",");
  const seen = new Set<MapLayer>();
  for (const part of parts) {
    const layer = part.trim() as MapLayer;
    if ((MAP_LAYERS as readonly string[]).includes(layer)) {
      seen.add(layer);
    }
  }
  return seen.size > 0 ? MAP_LAYERS.filter((layer) => seen.has(layer)) : [...MAP_LAYERS];
}

export function mapQueryString(query: MapQuery): string {
  const params = new URLSearchParams();
  if (query.level !== "world") {
    params.set("level", query.level);
  }
  if (query.parent) {
    params.set("parent", query.parent);
  }
  if (query.country) {
    params.set("country", query.country);
  }
  if (query.region) {
    params.set("region", query.region);
  }
  if (query.network) {
    params.set("network", query.network);
  }
  if (query.service) {
    params.set("service", query.service);
  }
  if (query.q) {
    params.set("q", query.q);
  }
  if (query.select) {
    params.set("select", query.select);
  }
  if (!sameLayers(query.layers, MAP_LAYERS)) {
    for (const layer of query.layers) {
      params.append("layers", layer);
    }
  }
  const encoded = params.toString();
  return encoded ? `?${encoded}` : "";
}

export function mapAggregatesQueryString(
  query: MapQuery,
  viewport?: MapViewport | null
): string {
  const params = new URLSearchParams();
  params.set("level", query.level);
  if (query.parent) {
    params.set("parent", query.parent);
  }
  if (query.service) {
    params.set("service", query.service);
  }
  if (query.q) {
    params.set("q", query.q);
  }
  for (const layer of query.layers) {
    params.append("layers", layer);
  }
  if (viewport) {
    params.set("west", String(viewport.west));
    params.set("south", String(viewport.south));
    params.set("east", String(viewport.east));
    params.set("north", String(viewport.north));
  }
  return params.toString();
}

function sameLayers(a: readonly string[], b: readonly string[]): boolean {
  if (a.length !== b.length) {
    return false;
  }
  return a.every((layer, index) => layer === b[index]);
}
