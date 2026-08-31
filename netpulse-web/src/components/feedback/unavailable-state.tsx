import { UnplugIcon } from "lucide-react";

import { cn } from "@/lib/utils";

type UnavailableStateProps = {
  title?: string;
  description?: string;
  className?: string;
};

export function UnavailableState({
  title = "Unavailable",
  description = "This capability depends on a backend that is not connected yet.",
  className,
}: UnavailableStateProps) {
  return (
    <div
      role="status"
      className={cn(
        "flex flex-col items-start gap-3 rounded-lg border border-border bg-card px-4 py-6",
        className
      )}
    >
      <UnplugIcon className="size-5 text-muted-foreground" aria-hidden="true" />
      <div className="space-y-1">
        <h2 className="text-sm font-medium text-foreground">{title}</h2>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
    </div>
  );
}
