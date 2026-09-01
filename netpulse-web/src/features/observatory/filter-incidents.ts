import type {
  OutageQuery,
  PublicIncidentRecord,
} from "@/domain/observatory";
import { DEFAULT_OUTAGE_QUERY, PAGE_SIZE } from "@/domain/observatory";

export type FilteredIncidents = {
  items: PublicIncidentRecord[];
  total: number;
  page: number;
  pageSize: number;
};

const SEVERITY_RANK: Record<string, number> = {
  critical: 4,
  high: 3,
  moderate: 2,
  informational: 1,
};

export function parseOutageQuery(
  searchParams: Record<string, string | string[] | undefined>
): OutageQuery {
  const first = (key: string): string => {
    const value = searchParams[key];
    if (Array.isArray(value)) {
      return value[0] ?? "";
    }
    return value?.trim() ?? "";
  };
  const page = Number.parseInt(first("page"), 10);
  return {
    service: first("service"),
    region: first("region"),
    network: first("network"),
    severity: first("severity"),
    status: first("status"),
    time: first("time") || DEFAULT_OUTAGE_QUERY.time,
    q: first("q"),
    sort: first("sort") || DEFAULT_OUTAGE_QUERY.sort,
    page: Number.isFinite(page) && page > 0 ? page : 1,
  };
}

export function filterIncidents(
  items: PublicIncidentRecord[],
  query: OutageQuery,
  now = Date.now()
): FilteredIncidents {
  const filtered = items.filter((item) => {
    if (query.service && !matchesToken(item.affectedServices, query.service) && item.scope !== query.service) {
      return false;
    }
    if (query.region && !matchesToken(item.regions, query.region)) {
      return false;
    }
    if (query.network && !matchesToken(item.networks, query.network)) {
      return false;
    }
    if (query.severity && item.severity !== query.severity) {
      return false;
    }
    if (query.status && item.status !== query.status) {
      return false;
    }
    if (query.q && !matchesSearch(item, query.q)) {
      return false;
    }
    if (!inTimeWindow(item.startedAt, query.time, now)) {
      return false;
    }
    return true;
  });

  filtered.sort((a, b) => compareIncidents(a, b, query.sort));

  const start = (query.page - 1) * PAGE_SIZE;
  return {
    items: filtered.slice(start, start + PAGE_SIZE),
    total: filtered.length,
    page: query.page,
    pageSize: PAGE_SIZE,
  };
}

function matchesToken(values: string[], want: string): boolean {
  return values.some((value) => value.toLowerCase() === want.toLowerCase());
}

function matchesSearch(item: PublicIncidentRecord, raw: string): boolean {
  const q = raw.toLowerCase();
  return [
    item.title,
    item.scope,
    item.status,
    item.severity,
    ...item.affectedServices,
    ...item.regions,
    ...item.networks,
  ]
    .join(" ")
    .toLowerCase()
    .includes(q);
}

function inTimeWindow(startedAt: string, window: string, now: number): boolean {
  if (!window || window === "all") {
    return true;
  }
  const started = Date.parse(startedAt);
  if (Number.isNaN(started)) {
    return false;
  }
  const day = 24 * 60 * 60 * 1000;
  if (window === "24h") {
    return started >= now - day;
  }
  if (window === "7d") {
    return started >= now - 7 * day;
  }
  if (window === "30d") {
    return started >= now - 30 * day;
  }
  return true;
}

function compareIncidents(
  a: PublicIncidentRecord,
  b: PublicIncidentRecord,
  sort: string
): number {
  switch (sort) {
    case "started_asc":
      return a.startedAt.localeCompare(b.startedAt);
    case "severity":
      return (SEVERITY_RANK[b.severity] ?? 0) - (SEVERITY_RANK[a.severity] ?? 0);
    case "status":
      return a.status.localeCompare(b.status);
    case "updated_desc":
      return b.lastUpdatedAt.localeCompare(a.lastUpdatedAt);
    default:
      return b.startedAt.localeCompare(a.startedAt);
  }
}

export function outageQueryString(query: OutageQuery, page = query.page): string {
  const params = new URLSearchParams();
  const entries: [string, string][] = [
    ["service", query.service],
    ["region", query.region],
    ["network", query.network],
    ["severity", query.severity],
    ["status", query.status],
    ["time", query.time === "all" ? "" : query.time],
    ["q", query.q],
    ["sort", query.sort === "started_desc" ? "" : query.sort],
    ["page", page > 1 ? String(page) : ""],
  ];
  for (const [key, value] of entries) {
    if (value) {
      params.set(key, value);
    }
  }
  const encoded = params.toString();
  return encoded ? `?${encoded}` : "";
}
