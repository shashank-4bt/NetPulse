import { OBSERVED_FAILURES_TITLE } from "@/domain/observatory";

export type AttributionInput = {
  evidenceClass: string;
  isolatedLayer: string | null;
  namedNetwork: string | null;
};

export function incidentTitle(input: AttributionInput): string {
  if (
    input.evidenceClass === "measured_fact" &&
    input.isolatedLayer === "isp" &&
    input.namedNetwork
  ) {
    return `Measured failures isolated to network ${input.namedNetwork}`;
  }
  if (
    input.evidenceClass === "measured_fact" &&
    input.isolatedLayer === "service"
  ) {
    return "Elevated connectivity failures observed toward the service";
  }
  return OBSERVED_FAILURES_TITLE;
}

export function isCausalOverclaim(title: string): boolean {
  const lower = title.toLowerCase();
  return (
    lower.includes("caused the outage") ||
    lower.includes("is down for everyone")
  );
}
