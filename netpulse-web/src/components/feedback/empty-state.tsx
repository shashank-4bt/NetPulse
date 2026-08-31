import type { ReactNode } from "react";
import { InboxIcon } from "lucide-react";

import { cn } from "@/lib/utils";

type EmptyStateProps = {
  title: string;
  description: string;
  className?: string;
  children?: ReactNode;
};

export function EmptyState({
  title,
  description,
  className,
  children,
}: EmptyStateProps) {
  return (
    <div
      className={cn(
        "flex flex-col items-start gap-3 rounded-lg border border-dashed border-border bg-muted/30 px-4 py-6",
        className
      )}
    >
      <InboxIcon className="size-5 text-muted-foreground" aria-hidden="true" />
      <div className="space-y-1">
        <h2 className="text-sm font-medium text-foreground">{title}</h2>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
      {children}
    </div>
  );
}
