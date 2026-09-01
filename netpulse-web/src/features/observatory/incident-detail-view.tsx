import Link from "next/link";

import { EvidenceItem } from "@/components/data/evidence-item";
import { Timeline } from "@/components/data/timeline";
import { EmptyState } from "@/components/feedback/empty-state";
import { ConfidenceBadge } from "@/components/status/confidence-badge";
import { LayerStatusBadge } from "@/components/status/layer-status-badge";
import { SeverityBadge } from "@/components/status/severity-badge";
import { Badge } from "@/components/ui/badge";
import type { Severity } from "@/domain/display";
import { SEVERITIES } from "@/domain/display";
import type { PublicIncidentRecord } from "@/domain/observatory";
import { INCIDENT_STAGE_LABELS, INCIDENT_STAGES } from "@/domain/observatory";
import { sampleCaption } from "@/features/observatory/empty-intelligence";
import { timelineFor } from "@/features/observatory/timeline";

type IncidentDetailViewProps = {
  incident: PublicIncidentRecord;
};

export function IncidentDetailView({ incident }: IncidentDetailViewProps) {
  const events = timelineFor(incident);
  const severity = isSeverity(incident.severity) ? incident.severity : null;

  return (
    <div className="space-y-10">
      <section className="flex flex-wrap items-center gap-2" aria-label="Incident status">
        {severity ? <SeverityBadge severity={severity} /> : (
          <Badge variant="outline">Severity not classified</Badge>
        )}
        <Badge variant="outline">
          {isStage(incident.status)
            ? INCIDENT_STAGE_LABELS[incident.status]
            : "Status not classified"}
        </Badge>
      </section>

      <dl className="grid gap-4 sm:grid-cols-2">
        <div>
          <dt className="text-sm font-medium">Started</dt>
          <dd className="mt-1 font-mono text-sm text-muted-foreground">
            {incident.startedAt || "Not recorded"}
          </dd>
        </div>
        <div>
          <dt className="text-sm font-medium">Last updated</dt>
          <dd className="mt-1 font-mono text-sm text-muted-foreground">
            {incident.lastUpdatedAt || "Not recorded"}
          </dd>
        </div>
      </dl>

      <TokenList title="Affected services" values={incident.affectedServices} hrefFor={(slug) => `/service/${slug}`} />
      <TokenList title="Regions" values={incident.regions} />
      <TokenList title="Networks" values={incident.networks} />

      <section aria-labelledby="sample-heading">
        <h2 id="sample-heading" className="text-lg font-semibold">
          Observed sample
        </h2>
        <p className="mt-2 text-sm text-muted-foreground">
          {sampleCaption(incident.sampleCount)}
          {incident.sampleRate ? ` Rate: ${incident.sampleRate}.` : ""}
        </p>
        <p className="mt-2 text-sm text-muted-foreground">
          Affected-user counts are omitted. NetPulse does not invent population impact.
        </p>
      </section>

      <section aria-labelledby="evidence-heading">
        <h2 id="evidence-heading" className="text-lg font-semibold">
          Evidence
        </h2>
        {incident.evidence.length === 0 ? (
          <EmptyState
            className="mt-3"
            title="No evidence recorded"
            description="Facts appear here only after probes persist observations."
          />
        ) : (
          <ul className="mt-3 grid gap-3 lg:grid-cols-2">
            {incident.evidence.map((item) => (
              <li key={item.id}>
                <EvidenceItem
                  evidenceClass={item.evidenceClass}
                  title={item.title}
                  body={item.body}
                />
              </li>
            ))}
          </ul>
        )}
      </section>

      <section aria-labelledby="hypotheses-heading">
        <h2 id="hypotheses-heading" className="text-lg font-semibold">
          Hypotheses
        </h2>
        {incident.hypotheses.length === 0 ? (
          <EmptyState
            className="mt-3"
            title="No hypotheses"
            description="Inferences stay empty until evidence can support them. They are never shown as facts."
          />
        ) : (
          <ul className="mt-3 space-y-3">
            {incident.hypotheses.map((item) => (
              <li
                key={item.id}
                className="rounded-lg border border-dashed border-border bg-card p-4"
              >
                <p className="text-sm font-medium">{item.title}</p>
                <p className="mt-2 text-sm text-muted-foreground">{item.body}</p>
              </li>
            ))}
          </ul>
        )}
      </section>

      <section aria-labelledby="confidence-heading">
        <h2 id="confidence-heading" className="text-lg font-semibold">
          Confidence
        </h2>
        <div className="mt-3 space-y-2">
          {incident.confidence.level ? (
            <ConfidenceBadge confidence={incident.confidence.level} />
          ) : (
            <LayerStatusBadge status="insufficient_evidence" />
          )}
          <p className="text-sm text-muted-foreground">{incident.confidence.caveat}</p>
        </div>
      </section>

      <section aria-labelledby="timeline-heading">
        <h2 id="timeline-heading" className="text-lg font-semibold">
          Timeline
        </h2>
        <p className="mt-2 text-sm text-muted-foreground">
          Detected, Investigating, Correlated, Identified, Mitigating, Resolved.
          Resolved stays unreached until independent recoveries exist.
        </p>
        <div className="mt-4">
          <Timeline
            events={events.map((event) => ({
              id: event.stage,
              timestampLabel: event.at ?? event.status.replace(/_/g, " "),
              title: event.label,
              detail: event.note,
            }))}
          />
        </div>
      </section>
    </div>
  );
}

function isSeverity(value: string): value is Severity {
  return (SEVERITIES as readonly string[]).includes(value);
}

function isStage(value: string): value is (typeof INCIDENT_STAGES)[number] {
  return (INCIDENT_STAGES as readonly string[]).includes(value);
}

function TokenList({
  title,
  values,
  hrefFor,
}: {
  title: string;
  values: string[];
  hrefFor?: (value: string) => string;
}) {
  const headingId = `${title.toLowerCase().replace(/\s+/g, "-")}-heading`;
  return (
    <section aria-labelledby={headingId}>
      <h2 id={headingId} className="text-lg font-semibold">
        {title}
      </h2>
      {values.length === 0 ? (
        <p className="mt-2 text-sm text-muted-foreground">None recorded.</p>
      ) : (
        <ul className="mt-2 flex flex-wrap gap-2">
          {values.map((value) => (
            <li key={value}>
              {hrefFor ? (
                <Link href={hrefFor(value)} className="text-sm hover:underline">
                  {value}
                </Link>
              ) : (
                <Badge variant="outline">{value}</Badge>
              )}
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
