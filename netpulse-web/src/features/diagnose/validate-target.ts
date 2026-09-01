const BLOCKED_HOSTS = new Set([
  "localhost",
  "localhost.localdomain",
  "metadata.google.internal",
]);

const IPV4 =
  /^(?:(?:25[0-5]|2[0-4]\d|1?\d?\d)\.){3}(?:25[0-5]|2[0-4]\d|1?\d?\d)$/;

export type TargetValidation =
  | { ok: true; hostname: string }
  | { ok: false; error: string };

function isPrivateOrLocalIp(ip: string): boolean {
  const parts = ip.split(".").map((part) => Number(part));
  if (parts.length !== 4 || parts.some((part) => Number.isNaN(part))) {
    return true;
  }
  const [a, b] = parts;
  if (a === undefined || b === undefined) {
    return true;
  }
  if (a === 10 || a === 127 || a === 0) {
    return true;
  }
  if (a === 169 && b === 254) {
    return true;
  }
  if (a === 192 && b === 168) {
    return true;
  }
  if (a === 172 && b >= 16 && b <= 31) {
    return true;
  }
  return false;
}

export function validateDiagnoseTarget(raw: string): TargetValidation {
  const trimmed = raw.trim();
  if (!trimmed) {
    return { ok: false, error: "Enter a hostname or URL." };
  }

  let hostname = trimmed;
  try {
    if (trimmed.includes("://")) {
      const url = new URL(trimmed);
      hostname = url.hostname;
    }
  } catch {
    return { ok: false, error: "Enter a valid hostname or URL." };
  }

  hostname = hostname.replace(/\.$/, "").toLowerCase();

  if (!hostname || hostname.includes(" ")) {
    return { ok: false, error: "Enter a valid hostname or URL." };
  }

  if (BLOCKED_HOSTS.has(hostname) || hostname.endsWith(".local")) {
    return {
      ok: false,
      error: "Local and internal hostnames cannot be probed.",
    };
  }

  if (IPV4.test(hostname) && isPrivateOrLocalIp(hostname)) {
    return {
      ok: false,
      error: "Private, loopback, and link-local addresses cannot be probed.",
    };
  }

  if (!hostname.includes(".") && !IPV4.test(hostname)) {
    return { ok: false, error: "Enter a fully qualified hostname." };
  }

  return { ok: true, hostname };
}
