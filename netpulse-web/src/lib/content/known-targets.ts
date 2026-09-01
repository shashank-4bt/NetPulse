import { SERVICE_CATALOG } from "@/lib/content/services";

export type KnownTarget = {
  slug: string;
  name: string;
  hostname: string;
};

const HOST_OVERRIDES: Record<string, string> = {
  "microsoft-365": "office.com",
};

export function getKnownTargets(): KnownTarget[] {
  return SERVICE_CATALOG.map((service) => ({
    slug: service.slug,
    name: service.name,
    hostname: HOST_OVERRIDES[service.slug] ?? `${service.slug}.com`,
  }));
}

export function resolveKnownTarget(raw: string): KnownTarget | undefined {
  const normalized = raw.trim().toLowerCase();
  return getKnownTargets().find((target) => {
    return (
      target.slug === normalized ||
      target.name.toLowerCase() === normalized ||
      target.hostname === normalized ||
      `www.${target.hostname}` === normalized
    );
  });
}
