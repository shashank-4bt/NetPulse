import {
  INCIDENT_STAGE_LABELS,
  INCIDENT_STAGES,
  type IncidentTimelineEvent,
  type PublicIncidentRecord,
} from "@/domain/observatory";
import { canMarkResolved } from "@/features/observatory/resolution";

export function timelineFor(
  incident: PublicIncidentRecord
): IncidentTimelineEvent[] {
  const identified = incident.evidence.some(
    (item) => item.evidenceClass === "measured_fact"
  );
  const resolution = canMarkResolved({
    recoverySampleCount: incident.sampleCount,
    identifiedCause: identified,
  });
  const currentIndex = INCIDENT_STAGES.indexOf(
    incident.status as (typeof INCIDENT_STAGES)[number]
  );

  return INCIDENT_STAGES.map((stage, index) => {
    const stored = incident.timeline.find((event) => event.stage === stage);
    if (stage === "resolved" && !resolution.ok) {
      return {
        stage,
        label: INCIDENT_STAGE_LABELS[stage],
        status: "not_reached",
        at: null,
        note: resolution.reason,
      };
    }
    if (stored) {
      return stored;
    }
    let status: IncidentTimelineEvent["status"] = "not_reached";
    if (currentIndex >= 0) {
      if (index < currentIndex) {
        status = "complete";
      } else if (index === currentIndex) {
        status = "current";
      }
    }
    return {
      stage,
      label: INCIDENT_STAGE_LABELS[stage],
      status,
      at: status === "current" ? incident.startedAt : null,
      note:
        status === "not_reached"
          ? "No evidence for this stage yet."
          : "Stage recorded from stored incident status, not from a user-path probe.",
    };
  });
}
