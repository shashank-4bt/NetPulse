import { AlertOctagonIcon, InfoIcon, ShieldAlertIcon } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import type { Severity } from "@/domain/display";
import { SEVERITY_LABEL } from "@/lib/design/taxonomy";
import { cn } from "@/lib/utils";

const SEVERITY_ICON = {
  informational: InfoIcon,
  moderate: ShieldAlertIcon,
  high: AlertOctagonIcon,
  critical: AlertOctagonIcon,
} as const;

const SEVERITY_CLASS: Record<Severity, string> = {
  informational:
    "border-border bg-muted text-muted-foreground",
  moderate:
    "border-[color-mix(in_oklch,var(--status-degraded),transparent_55%)] bg-[color-mix(in_oklch,var(--status-degraded),transparent_88%)] text-[var(--status-degraded)]",
  high:
    "border-[color-mix(in_oklch,var(--status-incident),transparent_50%)] bg-[color-mix(in_oklch,var(--status-incident),transparent_90%)] text-[var(--status-incident)]",
  critical:
    "border-[color-mix(in_oklch,var(--status-incident),transparent_30%)] bg-[color-mix(in_oklch,var(--status-incident),transparent_82%)] text-[var(--status-incident)] font-semibold",
};

type SeverityBadgeProps = {
  severity: Severity;
  className?: string;
};

export function SeverityBadge({ severity, className }: SeverityBadgeProps) {
  const Icon = SEVERITY_ICON[severity];
  const label = SEVERITY_LABEL[severity];

  return (
    <Badge
      variant="outline"
      className={cn("h-7 rounded-md px-2", SEVERITY_CLASS[severity], className)}
    >
      <Icon aria-hidden="true" />
      <span>
        <span className="sr-only">Severity: </span>
        {label}
      </span>
    </Badge>
  );
}
