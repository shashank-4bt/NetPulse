/**
 * Public health adapter.
 * Until workers and stores exist, every query is unavailable.
 * Callers must not invent operational or incident values.
 */

export type PublicDataState = "unavailable";

export type PublicHealthSnapshot = {
  state: PublicDataState;
  reason: string;
  incidents: [];
  regions: [];
  measuredAt: null;
};

export function getPublicHealthSnapshot(): PublicHealthSnapshot {
  return {
    state: "unavailable",
    reason:
      "Measurement workers, PostgreSQL, and ClickHouse are not connected. NetPulse will not display invented service health, outages, or map series.",
    incidents: [],
    regions: [],
    measuredAt: null,
  };
}

