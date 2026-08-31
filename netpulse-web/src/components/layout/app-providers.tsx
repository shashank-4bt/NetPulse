"use client";

import type { ReactNode } from "react";

import { TooltipProvider } from "@/components/ui/tooltip";
import { Toaster } from "@/components/ui/sonner";
import { ThemeProvider } from "@/components/layout/theme-provider";

export function AppProviders({ children }: { children: ReactNode }) {
  return (
    <ThemeProvider>
      <TooltipProvider>
        {children}
        <Toaster />
      </TooltipProvider>
    </ThemeProvider>
  );
}
