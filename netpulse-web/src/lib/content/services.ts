export type ServiceCatalogEntry = {
  slug: string;
  name: string;
  category: string;
  summary: string;
  layers: readonly string[];
};

/**
 * Editorial catalog of services NetPulse can diagnose.
 * This is not live health data. Pages must not imply current status.
 */
export const SERVICE_CATALOG: readonly ServiceCatalogEntry[] = [
  {
    slug: "google",
    name: "Google",
    category: "Search & identity",
    summary: "Search, accounts, and related Google properties.",
    layers: ["DNS", "Routing", "CDN", "TLS", "HTTP", "Service"],
  },
  {
    slug: "youtube",
    name: "YouTube",
    category: "Media",
    summary: "Video delivery and playback endpoints.",
    layers: ["DNS", "CDN", "TLS", "HTTP", "Service"],
  },
  {
    slug: "cloudflare",
    name: "Cloudflare",
    category: "Infrastructure",
    summary: "DNS, CDN, and edge connectivity used by many sites.",
    layers: ["DNS", "Connectivity", "Routing", "CDN", "TLS"],
  },
  {
    slug: "github",
    name: "GitHub",
    category: "Developer",
    summary: "Git hosting, APIs, and web application.",
    layers: ["DNS", "TLS", "HTTP", "Service"],
  },
  {
    slug: "microsoft-365",
    name: "Microsoft 365",
    category: "Productivity",
    summary: "Identity, mail, and collaboration endpoints.",
    layers: ["DNS", "ISP", "TLS", "HTTP", "Service"],
  },
  {
    slug: "slack",
    name: "Slack",
    category: "Collaboration",
    summary: "Messaging clients and real-time service endpoints.",
    layers: ["DNS", "CDN", "TLS", "HTTP", "Service"],
  },
  {
    slug: "aws",
    name: "AWS",
    category: "Cloud",
    summary: "Regional control-plane and commonly used public endpoints.",
    layers: ["DNS", "Routing", "TLS", "HTTP", "Service"],
  },
  {
    slug: "zoom",
    name: "Zoom",
    category: "Meetings",
    summary: "Meeting join paths and media edge connectivity.",
    layers: ["Connectivity", "ISP", "Routing", "TLS", "Service"],
  },
] as const;

export function getServiceBySlug(
  slug: string
): ServiceCatalogEntry | undefined {
  return SERVICE_CATALOG.find((service) => service.slug === slug);
}

export function getServiceSlugs(): string[] {
  return SERVICE_CATALOG.map((service) => service.slug);
}
