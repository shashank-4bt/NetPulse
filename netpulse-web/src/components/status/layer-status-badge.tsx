import {
  CircleAlertIcon,
  CircleCheckIcon,
  CircleHelpIcon,
  CircleMinusIcon,
  TriangleAlertIcon,
} from "lucide-react";

import { Badge } from "@/components/ui/badge";
import type { LayerStatus } from "@/domain/display";
import { LAYER_STATUS_LABEL } from "@/lib/design/taxonomy";
import { cn } from "@/lib/utils";

const STATUS_ICON = {
  not_measured: CircleMinusIcon,
  insufficient_evidence: CircleHelpIcon,
  healthy: CircleCheckIcon,
  degraded: TriangleAlertIcon,
  failed: CircleAlertIcon,
} as const;

type LayerStatusBadgeProps = {
  status: LayerStatus;
  className?: string;
};

export function LayerStatusBadge({ status, className }: LayerStatusBadgeProps) {
  const Icon = STATUS_ICON[status];

  return (
    <Badge variant="outline" className={cn("h-7 rounded-md px-2", className)}>
      <Icon aria-hidden="true" />
      <span>{LAYER_STATUS_LABEL[status]}</span>
    </Badge>
  );
}
