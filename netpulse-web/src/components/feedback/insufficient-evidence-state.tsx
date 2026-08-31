import { CircleHelpIcon } from "lucide-react";

import { INSUFFICIENT_EVIDENCE_LABEL } from "@/lib/design/taxonomy";
import { cn } from "@/lib/utils";

type InsufficientEvidenceStateProps = {
  description?: string;
  className?: string;
};

export function InsufficientEvidenceState({
  description = "Not enough measured observations are available to support a conclusion.",
  className,
}: InsufficientEvidenceStateProps) {
  return (
    <div
      role="status"
      className={cn(
        "flex flex-col items-start gap-3 rounded-lg border border-border bg-muted/20 px-4 py-6",
        className
      )}
    >
      <CircleHelpIcon className="size-5 text-muted-foreground" aria-hidden="true" />
      <div className="space-y-1">
        <h2 className="text-sm font-medium text-foreground">
          {INSUFFICIENT_EVIDENCE_LABEL}
        </h2>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
    </div>
  );
}
