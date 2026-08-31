import type { ReactNode } from "react";

import { EmptyState } from "@/components/feedback/empty-state";
import { ErrorState } from "@/components/feedback/error-state";
import { InsufficientEvidenceState } from "@/components/feedback/insufficient-evidence-state";
import { LoadingState } from "@/components/feedback/loading-state";
import { UnavailableState } from "@/components/feedback/unavailable-state";
import { cn } from "@/lib/utils";

type ChartContainerState =
  | "ready"
  | "loading"
  | "empty"
  | "error"
  | "unavailable"
  | "insufficient_evidence";

type ChartContainerProps = {
  title: string;
  description?: string;
  state?: ChartContainerState;
  children?: ReactNode;
  className?: string;
};

export function ChartContainer({
  title,
  description,
  state = "empty",
  children,
  className,
}: ChartContainerProps) {
  return (
    <section
      className={cn(
        "flex flex-col gap-3 rounded-lg border border-border bg-card p-4",
        className
      )}
      aria-labelledby={`${slug(title)}-chart-title`}
    >
      <header>
        <h2
          id={`${slug(title)}-chart-title`}
          className="text-sm font-medium text-foreground"
        >
          {title}
        </h2>
        {description ? (
          <p className="text-sm text-muted-foreground">{description}</p>
        ) : null}
      </header>
      <div className="min-h-40">
        {state === "ready" ? children : null}
        {state === "loading" ? <LoadingState label="Loading chart" /> : null}
        {state === "empty" ? (
          <EmptyState
            title="No chart data"
            description="A chart is rendered only when measured series are supplied. This container will not invent values."
          />
        ) : null}
        {state === "error" ? (
          <ErrorState
            title="Chart failed to load"
            description="The data source returned an error. No substitute series were generated."
          />
        ) : null}
        {state === "unavailable" ? (
          <UnavailableState description="The metrics API is not connected." />
        ) : null}
        {state === "insufficient_evidence" ? (
          <InsufficientEvidenceState description="There are not enough measured points to plot a reliable series." />
        ) : null}
      </div>
    </section>
  );
}

function slug(value: string): string {
  return value.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-|-$/g, "");
}
