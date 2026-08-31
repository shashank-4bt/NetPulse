import { CircleAlertIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

type ErrorStateProps = {
  title: string;
  description: string;
  className?: string;
  retryLabel?: string;
  onRetry?: () => void;
};

export function ErrorState({
  title,
  description,
  className,
  retryLabel = "Try again",
  onRetry,
}: ErrorStateProps) {
  return (
    <div
      role="alert"
      className={cn(
        "flex flex-col items-start gap-3 rounded-lg border border-destructive/30 bg-destructive/5 px-4 py-6",
        className
      )}
    >
      <CircleAlertIcon className="size-5 text-destructive" aria-hidden="true" />
      <div className="space-y-1">
        <h2 className="text-sm font-medium text-foreground">{title}</h2>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
      {onRetry ? (
        <Button type="button" variant="outline" onClick={onRetry}>
          {retryLabel}
        </Button>
      ) : null}
    </div>
  );
}
