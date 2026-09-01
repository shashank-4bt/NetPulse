import type { ReactNode } from "react";

import type { OutageQuery } from "@/domain/observatory";
import {
  INCIDENT_STAGE_LABELS,
  INCIDENT_STAGES,
  OUTAGE_SORTS,
  OUTAGE_TIME_WINDOWS,
} from "@/domain/observatory";
import { SEVERITIES } from "@/domain/display";
import { SEVERITY_LABEL } from "@/lib/design/taxonomy";
import type { ServiceCatalogEntry } from "@/lib/content/services";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

type OutageFiltersProps = {
  query: OutageQuery;
  services: readonly ServiceCatalogEntry[];
};

const FIELD =
  "h-11 min-h-11 w-full rounded-lg border border-input bg-transparent px-3 text-sm outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50";

export function OutageFilters({ query, services }: OutageFiltersProps) {
  return (
    <form method="get" action="/outages" className="space-y-4">
      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        <Field label="Search" htmlFor="q">
          <Input id="q" name="q" defaultValue={query.q} placeholder="Title, service, region, network" />
        </Field>
        <Field label="Service" htmlFor="service">
          <select id="service" name="service" defaultValue={query.service} className={FIELD}>
            <option value="">All services</option>
            {services.map((service) => (
              <option key={service.slug} value={service.slug}>
                {service.name}
              </option>
            ))}
          </select>
        </Field>
        <Field label="Region" htmlFor="region">
          <Input id="region" name="region" defaultValue={query.region} placeholder="e.g. eu-west" />
        </Field>
        <Field label="Network" htmlFor="network">
          <Input id="network" name="network" defaultValue={query.network} placeholder="e.g. AS64500" />
        </Field>
        <Field label="Severity" htmlFor="severity">
          <select id="severity" name="severity" defaultValue={query.severity} className={FIELD}>
            <option value="">All severities</option>
            {SEVERITIES.map((severity) => (
              <option key={severity} value={severity}>
                {SEVERITY_LABEL[severity]}
              </option>
            ))}
          </select>
        </Field>
        <Field label="Status" htmlFor="status">
          <select id="status" name="status" defaultValue={query.status} className={FIELD}>
            <option value="">All statuses</option>
            {INCIDENT_STAGES.map((stage) => (
              <option key={stage} value={stage}>
                {INCIDENT_STAGE_LABELS[stage]}
              </option>
            ))}
          </select>
        </Field>
        <Field label="Time" htmlFor="time">
          <select id="time" name="time" defaultValue={query.time} className={FIELD}>
            {OUTAGE_TIME_WINDOWS.map((window) => (
              <option key={window} value={window}>
                {window === "all" ? "All time" : window}
              </option>
            ))}
          </select>
        </Field>
        <Field label="Sort" htmlFor="sort">
          <select id="sort" name="sort" defaultValue={query.sort} className={FIELD}>
            {OUTAGE_SORTS.map((sort) => (
              <option key={sort} value={sort}>
                {sortLabel(sort)}
              </option>
            ))}
          </select>
        </Field>
      </div>
      <Button type="submit">Apply filters</Button>
    </form>
  );
}

function Field({
  label,
  htmlFor,
  children,
}: {
  label: string;
  htmlFor: string;
  children: ReactNode;
}) {
  return (
    <div className="space-y-1.5">
      <label htmlFor={htmlFor} className="text-sm font-medium">
        {label}
      </label>
      {children}
    </div>
  );
}

function sortLabel(value: string): string {
  switch (value) {
    case "started_asc":
      return "Oldest first";
    case "updated_desc":
      return "Recently updated";
    case "severity":
      return "Severity";
    case "status":
      return "Status";
    default:
      return "Newest first";
  }
}
