import type { ReactNode } from "react";

import { cn } from "@/lib/utils";

type SectionContainerProps = {
  children: ReactNode;
  className?: string;
  id?: string;
  labelledBy?: string;
};

export function SectionContainer({
  children,
  className,
  id,
  labelledBy,
}: SectionContainerProps) {
  return (
    <section
      id={id}
      aria-labelledby={labelledBy}
      className={cn("py-8 md:py-10", className)}
    >
      {children}
    </section>
  );
}
