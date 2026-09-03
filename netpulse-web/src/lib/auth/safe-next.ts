const FALLBACK = "/dashboard";

const ALLOWED_PREFIXES = [
  "/dashboard",
  "/account",
  "/developers",
  "/business",
  "/admin",
  "/diagnose",
  "/reports",
  "/map",
  "/services",
  "/status",
  "/outages",
  "/trust",
  "/privacy",
  "/security",
  "/how-it-works",
] as const;

export function safeNext(value: string, fallback = FALLBACK): string {
  if (!value) {
    return fallback;
  }
  const trimmed = value.trim();
  if (
    trimmed.includes("\\") ||
    trimmed.includes("@") ||
    trimmed.includes("://") ||
    !trimmed.startsWith("/") ||
    trimmed.startsWith("//")
  ) {
    return fallback;
  }

  let parsed: URL;
  try {
    parsed = new URL(trimmed, "https://netpulse.invalid");
  } catch {
    return fallback;
  }
  if (parsed.origin !== "https://netpulse.invalid") {
    return fallback;
  }
  if (parsed.username || parsed.password) {
    return fallback;
  }
  if (parsed.pathname.startsWith("//")) {
    return fallback;
  }

  const allowed = ALLOWED_PREFIXES.some(
    (prefix) => parsed.pathname === prefix || parsed.pathname.startsWith(`${prefix}/`)
  );
  if (!allowed) {
    return fallback;
  }

  return `${parsed.pathname}${parsed.search}${parsed.hash}`;
}
