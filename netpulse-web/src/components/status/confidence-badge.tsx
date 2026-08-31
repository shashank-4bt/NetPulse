import { GaugeIcon } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import type { ConfidenceLevel } from "@/domain/display";
import { CONFIDENCE_LABEL } from "@/lib/design/taxonomy";
import { cn } from "@/lib/utils";

const CONFIDENCE_MARK: Record<ConfidenceLevel, string> = {
  low: "●○○○",
  medium: "●●○○",
  high: "●●●○",
  very_high: "●●●●",
};

type ConfidenceBadgeProps = {
  confidence: ConfidenceLevel;
  className?: string;
};

export function ConfidenceBadge({
  confidence,
  className,
}: ConfidenceBadgeProps) {
  const label = CONFIDENCE_LABEL[confidence];

  return (
    <Badge variant="outline" className={cn("h-7 rounded-md px-2", className)}>
      <GaugeIcon aria-hidden="true" />
      <span aria-hidden="true" className="font-mono text-[0.65rem] tracking-tight">
        {CONFIDENCE_MARK[confidence]}
      </span>
      <span>
        <span className="sr-only">Confidence: </span>
        {label}
      </span>
    </Badge>
  );
}
