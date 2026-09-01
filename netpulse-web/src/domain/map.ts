export const MAP_LEVELS = [
  "world",
  "country",
  "region",
  "network",
  "service",
] as const;

export type MapLevel = (typeof MAP_LEVELS)[number];

export const MAP_LAYERS = [
  "global",
  "regional",
  "network",
  "service",
  "incidents",
] as const;

export type MapLayer = (typeof MAP_LAYERS)[number];

export const MAP_CELL_LIMIT = 250;
export const MAP_INCIDENT_LIMIT = 50;

export const MAP_LAYER_LABELS: Record<MapLayer, string> = {
  global: "Global Health",
  regional: "Regional Health",
  network: "Network Health",
  service: "Service Health",
  incidents: "Incidents",
};

export const MAP_LEVEL_LABELS: Record<MapLevel, string> = {
  world: "World",
  country: "Country",
  region: "Region",
  network: "Network/ASN",
  service: "Service",
};

export const MAP_STATUSES = [
  "not_measured",
  "insufficient_evidence",
  "healthy",
  "degraded",
  "failed",
] as const;

export type MapCellStatus = (typeof MAP_STATUSES)[number];

export type MapCell = {
  id: string;
  level: MapLevel | string;
  label: string;
  parentId: string | null;
  lon: number | null;
  lat: number | null;
  sampleCount: number;
  status: string;
  summary: string;
  layer: MapLayer | string;
  childCount: number;
};

export type MapIncidentRef = {
  id: string;
  title: string;
  coarseRegion: string;
  sampleCount: number;
};

export type MapAggregates = {
  level: MapLevel | string;
  parentId: string | null;
  cells: MapCell[];
  incidentRefs: MapIncidentRef[];
  totalSamples: number;
  limit: number;
  truncated: boolean;
  precision: "none" | "degree" | "country" | "region" | string;
  reason: string;
  hasCoordinates: boolean;
};

export type MapViewport = {
  west: number;
  south: number;
  east: number;
  north: number;
};

export type MapQuery = {
  level: MapLevel;
  parent: string;
  country: string;
  region: string;
  network: string;
  service: string;
  q: string;
  layers: MapLayer[];
  select: string;
};

export const DEFAULT_MAP_QUERY: MapQuery = {
  level: "world",
  parent: "",
  country: "",
  region: "",
  network: "",
  service: "",
  q: "",
  layers: [...MAP_LAYERS],
  select: "",
};

export function emptyMapAggregates(reason: string): MapAggregates {
  return {
    level: "world",
    parentId: null,
    cells: [],
    incidentRefs: [],
    totalSamples: 0,
    limit: MAP_CELL_LIMIT,
    truncated: false,
    precision: "none",
    reason,
    hasCoordinates: false,
  };
}
