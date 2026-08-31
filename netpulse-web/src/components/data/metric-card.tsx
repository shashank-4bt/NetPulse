import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { INSUFFICIENT_EVIDENCE_LABEL } from "@/lib/design/taxonomy";
import { cn } from "@/lib/utils";

type MetricCardProps = {
  title: string;
  description?: string;
  /** Null means the value is not available — never invent a number. */
  value: string | number | null;
  unit?: string;
  caption?: string;
  className?: string;
};

export function MetricCard({
  title,
  description,
  value,
  unit,
  caption,
  className,
}: MetricCardProps) {
  const missing = value === null;

  return (
    <Card className={cn(className)}>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        {description ? <CardDescription>{description}</CardDescription> : null}
      </CardHeader>
      <CardContent>
        {missing ? (
          <p className="text-sm text-muted-foreground" role="status">
            {INSUFFICIENT_EVIDENCE_LABEL}
          </p>
        ) : (
          <p className="font-mono text-2xl font-semibold tabular-nums tracking-tight">
            {value}
            {unit ? (
              <span className="ml-1 text-sm font-normal text-muted-foreground">
                {unit}
              </span>
            ) : null}
          </p>
        )}
        {caption ? (
          <p className="mt-2 text-xs text-muted-foreground">{caption}</p>
        ) : null}
        {missing ? (
          <Badge variant="outline" className="mt-3">
            Not measured
          </Badge>
        ) : null}
      </CardContent>
    </Card>
  );
}
