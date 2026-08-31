"use client";

import { motion, useReducedMotion } from "framer-motion";

import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

type LoadingStateProps = {
  label?: string;
  className?: string;
};

export function LoadingState({
  label = "Loading",
  className,
}: LoadingStateProps) {
  const reduceMotion = useReducedMotion();

  return (
    <div
      role="status"
      aria-live="polite"
      aria-busy="true"
      className={cn("flex flex-col gap-3", className)}
    >
      <p className="text-sm text-muted-foreground">{label}…</p>
      <motion.div
        className="flex flex-col gap-2"
        animate={reduceMotion ? undefined : { opacity: [0.55, 1, 0.55] }}
        transition={
          reduceMotion
            ? undefined
            : { duration: 1.6, repeat: Infinity, ease: "easeInOut" }
        }
      >
        <Skeleton className="h-4 w-2/3" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-5/6" />
      </motion.div>
    </div>
  );
}
