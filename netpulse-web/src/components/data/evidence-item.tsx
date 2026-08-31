import {
  FlaskConicalIcon,
  LightbulbIcon,
  RulerIcon,
} from "lucide-react";

import { Badge } from "@/components/ui/badge";
import type { EvidenceClass } from "@/domain/display";
import { EVIDENCE_CLASS_LABEL } from "@/lib/design/taxonomy";
import { cn } from "@/lib/utils";

const CLASS_ICON = {
  measured_fact: RulerIcon,
  inferred_hypothesis: FlaskConicalIcon,
  recommendation: LightbulbIcon,
} as const;

type EvidenceItemProps = {
  evidenceClass: EvidenceClass;
  title: string;
  body: string;
  source?: string;
  className?: string;
};

export function EvidenceItem({
  evidenceClass,
  title,
  body,
  source,
  className,
}: EvidenceItemProps) {
  const Icon = CLASS_ICON[evidenceClass];
  const classLabel = EVIDENCE_CLASS_LABEL[evidenceClass];

  return (
    <article
      className={cn(
        "flex flex-col gap-2 rounded-lg border border-border bg-card p-4",
        evidenceClass === "inferred_hypothesis" && "border-dashed",
        className
      )}
    >
      <div className="flex flex-wrap items-center gap-2">
        <Badge variant="outline" className="h-7 rounded-md">
          <Icon aria-hidden="true" />
          <span>{classLabel}</span>
        </Badge>
      </div>
      <h3 className="text-sm font-medium text-foreground">{title}</h3>
      <p className="text-sm text-muted-foreground">{body}</p>
      {source ? (
        <p className="font-mono text-xs text-muted-foreground">{source}</p>
      ) : null}
    </article>
  );
}
