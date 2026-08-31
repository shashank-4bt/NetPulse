import {
  AlertTriangleIcon,
  CheckCircle2Icon,
  HelpCircleIcon,
  SearchIcon,
  SirenIcon,
} from "lucide-react";

import { Badge } from "@/components/ui/badge";
import type { OperationalStatus } from "@/domain/display";
import { STATUS_LABEL } from "@/lib/design/taxonomy";
import { cn } from "@/lib/utils";

const STATUS_ICON = {
  operational: CheckCircle2Icon,
  degraded: AlertTriangleIcon,
  investigating: SearchIcon,
  major_incident: SirenIcon,
  unknown: HelpCircleIcon,
} as const;

const STATUS_CLASS: Record<OperationalStatus, string> = {
  operational:
    "border-[color-mix(in_oklch,var(--status-operational),transparent_55%)] bg-[color-mix(in_oklch,var(--status-operational),transparent_88%)] text-[var(--status-operational)]",
  degraded:
    "border-[color-mix(in_oklch,var(--status-degraded),transparent_55%)] bg-[color-mix(in_oklch,var(--status-degraded),transparent_88%)] text-[var(--status-degraded)]",
  investigating:
    "border-[color-mix(in_oklch,var(--status-investigating),transparent_55%)] bg-[color-mix(in_oklch,var(--status-investigating),transparent_88%)] text-[var(--status-investigating)]",
  major_incident:
    "border-[color-mix(in_oklch,var(--status-incident),transparent_55%)] bg-[color-mix(in_oklch,var(--status-incident),transparent_88%)] text-[var(--status-incident)]",
  unknown:
    "border-[color-mix(in_oklch,var(--status-unknown),transparent_55%)] bg-[color-mix(in_oklch,var(--status-unknown),transparent_88%)] text-[var(--status-unknown)]",
};

type StatusBadgeProps = {
  status: OperationalStatus;
  className?: string;
};

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const Icon = STATUS_ICON[status];
  const label = STATUS_LABEL[status];

  return (
    <Badge
      variant="outline"
      className={cn("h-7 rounded-md px-2", STATUS_CLASS[status], className)}
    >
      <Icon aria-hidden="true" />
      <span>{label}</span>
    </Badge>
  );
}
