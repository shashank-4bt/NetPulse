"use client";

import { MonitorIcon, MoonIcon, SunIcon } from "lucide-react";
import { useTheme } from "next-themes";
import { useSyncExternalStore } from "react";

import { Button } from "@/components/ui/button";

const ORDER = ["light", "dark", "system"] as const;

function subscribe() {
  return () => undefined;
}

function useIsClient() {
  return useSyncExternalStore(subscribe, () => true, () => false);
}

export function ThemeToggle() {
  const { theme, setTheme } = useTheme();
  const isClient = useIsClient();

  const current = ORDER.includes(theme as (typeof ORDER)[number])
    ? (theme as (typeof ORDER)[number])
    : "system";

  function cycleTheme() {
    const index = ORDER.indexOf(current);
    const next = ORDER[(index + 1) % ORDER.length];
    if (next) {
      setTheme(next);
    }
  }

  const label =
    current === "light"
      ? "Light theme"
      : current === "dark"
        ? "Dark theme"
        : "System theme";

  return (
    <Button
      type="button"
      variant="outline"
      size="icon"
      className="size-11"
      onClick={cycleTheme}
      aria-label={`Color theme: ${label}. Activate to change.`}
    >
      {isClient && current === "dark" ? (
        <MoonIcon aria-hidden="true" />
      ) : isClient && current === "light" ? (
        <SunIcon aria-hidden="true" />
      ) : (
        <MonitorIcon aria-hidden="true" />
      )}
    </Button>
  );
}
