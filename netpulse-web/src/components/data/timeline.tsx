import { cn } from "@/lib/utils";

export type TimelineEvent = {
  id: string;
  timestampLabel: string;
  title: string;
  detail: string;
};

type TimelineProps = {
  events: TimelineEvent[];
  className?: string;
};

export function Timeline({ events, className }: TimelineProps) {
  if (events.length === 0) {
    return (
      <p className="text-sm text-muted-foreground" role="status">
        No timeline events.
      </p>
    );
  }

  return (
    <ol className={cn("relative space-y-4 border-l border-border pl-4", className)}>
      {events.map((event) => (
        <li key={event.id} className="relative">
          <span
            aria-hidden="true"
            className="absolute top-1.5 -left-[1.3rem] size-2.5 rounded-full border border-border bg-background"
          />
          <p className="font-mono text-xs text-muted-foreground">
            {event.timestampLabel}
          </p>
          <h3 className="text-sm font-medium text-foreground">{event.title}</h3>
          <p className="text-sm text-muted-foreground">{event.detail}</p>
        </li>
      ))}
    </ol>
  );
}
