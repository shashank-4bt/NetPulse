/**
 * Public runtime configuration. Server secrets must never be added here.
 * When the API is absent, callers must surface "unavailable" — not mock health.
 */
export const publicConfig = {
  appName: "NetPulse",
  appTagline: "Internet Health Intelligence Platform",
} as const;
