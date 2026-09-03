import { resolveKnownTarget } from "@/lib/content/known-targets";

const MAX_INPUT_LENGTH = 2048;
const MAX_HOSTNAME_LENGTH = 253;

const BLOCKED_HOSTS = new Set([
  "localhost",
  "localhost.localdomain",
  "metadata.google.internal",
]);

const UNSAFE_PROTOCOLS = new Set([
  "javascript:",
  "data:",
  "file:",
  "ftp:",
  "ws:",
  "wss:",
]);

const IPV4 =
  /^(?:(?:25[0-5]|2[0-4]\d|1?\d?\d)\.){3}(?:25[0-5]|2[0-4]\d|1?\d?\d)$/;

const HOSTNAME_LABEL = /^(?!-)[a-z0-9-]{1,63}(?<!-)$/i;

export type TargetKind = "domain" | "url" | "known_service";

export type TargetValidation =
  | {
      ok: true;
      hostname: string;
      kind: TargetKind;
      serviceSlug: string | null;
      raw: string;
    }
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
  if (a === 100 && b >= 64 && b <= 127) {
    return true;
  }
  if (a === 168 && b === 63 && parts[2] === 129 && parts[3] === 16) {
    return true;
  }
  if (a >= 240) {
    return true;
  }
  return false;
}

function looksLikeIpv6(value: string): boolean {
  return value.includes(":");
}

function isBlockedIpv6(value: string): boolean {
  const lowered = value.toLowerCase();
  return (
    lowered === "::1" ||
    lowered.startsWith("fe80:") ||
    lowered.startsWith("fc") ||
    lowered.startsWith("fd") ||
    lowered.startsWith("::ffff:127.") ||
    lowered.startsWith("::ffff:10.") ||
    lowered.startsWith("::ffff:192.168.") ||
    lowered.startsWith("::ffff:169.254.") ||
    lowered.startsWith("::ffff:168.63.129.16")
  );
}

function isValidHostname(hostname: string): boolean {
  if (IPV4.test(hostname)) {
    return !isPrivateOrLocalIp(hostname);
  }
  const labels = hostname.split(".");
  if (labels.length < 2) {
    return false;
  }
  return labels.every((label) => HOSTNAME_LABEL.test(label));
}

export function validateDiagnoseTarget(raw: string): TargetValidation {
  const trimmed = raw.trim();
  if (!trimmed) {
    return { ok: false, error: "Enter a hostname, URL, or known service." };
  }
  if (trimmed.length > MAX_INPUT_LENGTH) {
    return { ok: false, error: "Input is too long." };
  }
  if (/[<>'"\\]/.test(trimmed)) {
    return { ok: false, error: "Input contains unsafe characters." };
  }

  const known = resolveKnownTarget(trimmed);
  if (known) {
    return {
      ok: true,
      hostname: known.hostname,
      kind: "known_service",
      serviceSlug: known.slug,
      raw: trimmed,
    };
  }

  const lowered = trimmed.toLowerCase();
  for (const protocol of UNSAFE_PROTOCOLS) {
    if (lowered.startsWith(protocol)) {
      return { ok: false, error: "Only http and https URLs are accepted." };
    }
  }

  let hostname = trimmed;
  let kind: TargetKind = "domain";

  if (trimmed.includes("://")) {
    let url: URL;
    try {
      url = new URL(trimmed);
    } catch {
      return { ok: false, error: "Enter a valid hostname or URL." };
    }
    if (url.protocol !== "http:" && url.protocol !== "https:") {
      return { ok: false, error: "Only http and https URLs are accepted." };
    }
    if (url.username || url.password) {
      return { ok: false, error: "URLs with credentials are not accepted." };
    }
    hostname = url.hostname;
    kind = "url";
  }

  hostname = hostname.replace(/^\[|\]$/g, "").replace(/\.$/, "").toLowerCase();

  if (!hostname || hostname.includes(" ") || hostname.includes("/")) {
    return { ok: false, error: "Enter a valid hostname or URL." };
  }
  if (hostname.length > MAX_HOSTNAME_LENGTH) {
    return { ok: false, error: "Hostname is too long." };
  }
  if (BLOCKED_HOSTS.has(hostname) || hostname.endsWith(".local") || hostname.endsWith(".internal") || hostname.endsWith(".localhost") || hostname === "metadata") {
    return {
      ok: false,
      error: "Local and internal hostnames cannot be probed.",
    };
  }
  if (looksLikeIpv6(hostname) && isBlockedIpv6(hostname)) {
    return {
      ok: false,
      error: "Private, loopback, and link-local addresses cannot be probed.",
    };
  }
  if (IPV4.test(hostname) && isPrivateOrLocalIp(hostname)) {
    return {
      ok: false,
      error: "Private, loopback, and link-local addresses cannot be probed.",
    };
  }
  if (!isValidHostname(hostname) && !IPV4.test(hostname)) {
    return { ok: false, error: "Enter a fully qualified hostname." };
  }

  const matchedKnown = resolveKnownTarget(hostname);
  return {
    ok: true,
    hostname,
    kind: matchedKnown ? "known_service" : kind,
    serviceSlug: matchedKnown?.slug ?? null,
    raw: trimmed,
  };
}
