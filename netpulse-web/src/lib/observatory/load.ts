import type { ServiceCatalogEntry } from "@/lib/content/services";
import {
  getServiceBySlug,
  SERVICE_CATALOG,
} from "@/lib/content/services";
import type {
  OutageQuery,
  PublicIncidentRecord,
  ServiceIntelligence,
} from "@/domain/observatory";
import { emptyServiceIntelligence } from "@/features/observatory/empty-intelligence";
import { parseOutageQuery } from "@/features/observatory/filter-incidents";
import {
  getBackendIncident,
  getBackendIncidents,
  getBackendService,
  getBackendServices,
  isApiConfigured,
} from "@/lib/api/backend";
import { parseReportId } from "@/lib/reports/id";

export type ObservatoryState = "unavailable" | "empty" | "ready" | "error";

export type LoadedService = {
  catalog: ServiceCatalogEntry;
  intelligence: ServiceIntelligence;
  incidents: PublicIncidentRecord[];
  state: ObservatoryState;
  reason: string;
};

export type LoadedOutages = {
  query: OutageQuery;
  items: PublicIncidentRecord[];
  total: number;
  state: ObservatoryState;
  reason: string;
};

export type LoadedIncident = {
  incident: PublicIncidentRecord | null;
  missing: boolean;
  state: ObservatoryState;
  reason: string;
};

export async function loadServiceCatalog(): Promise<{
  services: ServiceCatalogEntry[];
  state: ObservatoryState;
  reason: string;
}> {
  if (!isApiConfigured()) {
    return {
      services: [...SERVICE_CATALOG],
      state: "unavailable",
      reason:
        "NETPULSE_API_BASE_URL is not set. Catalog entries are editorial and are not live health.",
    };
  }
  const result = await getBackendServices();
  if (!result.ok) {
    return {
      services: [...SERVICE_CATALOG],
      state: "error",
      reason: result.message,
    };
  }
  return {
    services: result.services.map((item) => ({
      slug: item.slug,
      name: item.name,
      category: item.category,
      summary: item.summary,
      layers: item.layers,
    })),
    state: "ready",
    reason: "Editorial catalog from the API. Current state still requires measured series.",
  };
}

export async function loadServicePage(slug: string): Promise<LoadedService | null> {
  const local = getServiceBySlug(slug);
  if (!isApiConfigured()) {
    if (!local) {
      return null;
    }
    return {
      catalog: local,
      intelligence: emptyServiceIntelligence(),
      incidents: [],
      state: "unavailable",
      reason:
        "The diagnose API is not configured. This page will not invent availability, latency, or incidents.",
    };
  }

  const [service, incidents] = await Promise.all([
    getBackendService(slug),
    getBackendIncidents({ service: slug, pageSize: 20 }),
  ]);

  if (!service.ok) {
    if (service.code === "not_found") {
      return null;
    }
    if (!local) {
      return null;
    }
    return {
      catalog: local,
      intelligence: emptyServiceIntelligence(),
      incidents: [],
      state: "error",
      reason: service.message,
    };
  }

  const catalog: ServiceCatalogEntry = {
    slug: service.service.slug,
    name: service.service.name,
    category: service.service.category,
    summary: service.service.summary,
    layers: service.service.layers,
  };
  const list = incidents.ok ? incidents.incidents : [];
  return {
    catalog,
    intelligence: service.intelligence,
    incidents: list,
    state: incidents.ok ? (list.length === 0 ? "empty" : "ready") : "error",
    reason: incidents.ok
      ? list.length === 0
        ? "No stored incidents for this service. That is not a claim that the service is healthy."
        : "Incidents listed here are stored records, not population impact."
      : incidents.message,
  };
}

export async function loadOutages(
  searchParams: Record<string, string | string[] | undefined>
): Promise<LoadedOutages> {
  const query = parseOutageQuery(searchParams);
  if (!isApiConfigured()) {
    return {
      query,
      items: [],
      total: 0,
      state: "unavailable",
      reason:
        "NETPULSE_API_BASE_URL is not set. The outage center will not invent incidents.",
    };
  }
  const result = await getBackendIncidents({
    service: query.service || undefined,
    region: query.region || undefined,
    network: query.network || undefined,
    severity: query.severity || undefined,
    status: query.status || undefined,
    time: query.time === "all" ? undefined : query.time,
    q: query.q || undefined,
    sort: query.sort,
    page: query.page,
    pageSize: 20,
  });
  if (!result.ok) {
    return {
      query,
      items: [],
      total: 0,
      state: "error",
      reason: result.message,
    };
  }
  return {
    query,
    items: result.incidents,
    total: result.total,
    state: result.total === 0 ? "empty" : "ready",
    reason:
      result.total === 0
        ? "The API store contains no matching incidents. An empty feed is not a healthy-internet claim."
        : "Rows are stored incident records. Sample counts are observed, not population impact.",
  };
}

export async function loadIncident(id: string): Promise<LoadedIncident> {
  const parsed = parseReportId(id);
  if (!parsed) {
    return {
      incident: null,
      missing: true,
      state: "empty",
      reason: "The incident id is not a valid identifier.",
    };
  }
  if (!isApiConfigured()) {
    return {
      incident: null,
      missing: false,
      state: "unavailable",
      reason:
        "NETPULSE_API_BASE_URL is not set. NetPulse will not invent an incident document.",
    };
  }
  const result = await getBackendIncident(parsed);
  if (!result.ok) {
    if (result.code === "not_found") {
      return {
        incident: null,
        missing: true,
        state: "empty",
        reason: "No stored incident has this id.",
      };
    }
    return {
      incident: null,
      missing: false,
      state: "error",
      reason: result.message,
    };
  }
  return {
    incident: result.incident,
    missing: false,
    state: "ready",
    reason: "This document is a stored incident. It is not a user-path diagnosis.",
  };
}
