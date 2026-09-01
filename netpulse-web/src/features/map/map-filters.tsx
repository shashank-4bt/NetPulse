import type { ReactNode } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  MAP_LAYER_LABELS,
  MAP_LAYERS,
  MAP_LEVEL_LABELS,
  MAP_LEVELS,
  type MapQuery,
} from "@/domain/map";
import type { ServiceCatalogEntry } from "@/lib/content/services";

type MapFiltersProps = {
  query: MapQuery;
  services: readonly ServiceCatalogEntry[];
};

const FIELD =
  "h-11 min-h-11 w-full rounded-lg border border-input bg-transparent px-3 text-sm outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50";

export function MapFilters({ query, services }: MapFiltersProps) {
  return (
    <form method="get" action="/map" className="space-y-4">
      <input type="hidden" name="parent" value={query.parent} />
      <input type="hidden" name="country" value={query.country} />
      <input type="hidden" name="region" value={query.region} />
      <input type="hidden" name="network" value={query.network} />
      <input type="hidden" name="select" value={query.select} />
      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        <Field label="Search" htmlFor="q">
          <Input
            id="q"
            name="q"
            defaultValue={query.q}
            placeholder="Region, network, service"
          />
        </Field>
        <Field label="Hierarchy" htmlFor="level">
          <select id="level" name="level" defaultValue={query.level} className={FIELD}>
            {MAP_LEVELS.map((level) => (
              <option key={level} value={level}>
                {MAP_LEVEL_LABELS[level]}
              </option>
            ))}
          </select>
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
      </div>
      <fieldset className="space-y-2">
        <legend className="text-sm font-medium">Layers</legend>
        <div className="flex flex-wrap gap-3">
          {MAP_LAYERS.map((layer) => (
            <label key={layer} className="inline-flex min-h-11 items-center gap-2 text-sm">
              <input
                type="checkbox"
                name="layers"
                value={layer}
                defaultChecked={query.layers.includes(layer)}
                className="size-4 rounded border-input"
              />
              {MAP_LAYER_LABELS[layer]}
            </label>
          ))}
        </div>
      </fieldset>
      <div className="flex flex-wrap gap-3">
        <Button type="submit">Apply filters</Button>
        <Button variant="outline" nativeButton={false} render={<a href="/map" />}>
          Reset to world
        </Button>
      </div>
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
