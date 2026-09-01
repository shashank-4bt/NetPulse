import { getBackendIncidents, isApiConfigured } from "@/lib/api/backend";

export type PublicDataState = "unavailable" | "empty" | "ready";

export type PublicIncident = {
  id: string;
  title: string;
  scope: string;
  startedAt: string;
};

export type PublicHealthSnapshot = {
  state: PublicDataState;
  reason: string;
  incidents: PublicIncident[];
  regions: [];
  measuredAt: null;
};

export async function getPublicHealthSnapshot(): Promise<PublicHealthSnapshot> {
  if (!isApiConfigured()) {
    return {
      state: "unavailable",
      reason:
        "NETPULSE_API_BASE_URL is not set. NetPulse will not display invented service health, outages, or map series.",
      incidents: [],
      regions: [],
      measuredAt: null,
    };
  }

  const result = await getBackendIncidents();
  if (!result.ok) {
    return {
      state: "unavailable",
      reason: result.message,
      incidents: [],
      regions: [],
      measuredAt: null,
    };
  }

  return {
    state: result.incidents.length === 0 ? "empty" : "ready",
    reason:
      result.incidents.length === 0
        ? "The API store contains no incidents. That is not a claim that the internet is healthy."
        : "Incidents are listed only as stored records. Regional health is still unmeasured.",
    incidents: result.incidents,
    regions: [],
    measuredAt: null,
  };
}
