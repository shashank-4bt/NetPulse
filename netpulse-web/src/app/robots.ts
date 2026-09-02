import type { MetadataRoute } from "next";

export default function robots(): MetadataRoute.Robots {
  const base = process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000";

  return {
    rules: {
      userAgent: "*",
      allow: "/",
      disallow: [
        "/design-system",
        "/reports",
        "/api",
        "/dashboard",
        "/account",
        "/share",
        "/developers/dashboard",
        "/developers/monitors",
        "/developers/incidents",
        "/developers/api",
        "/developers/webhooks",
        "/developers/usage",
        "/developers/sla",
        "/business/dashboard",
        "/business/devices",
        "/business/networks",
        "/business/services",
        "/business/incidents",
        "/business/analytics",
        "/business/reports",
        "/business/team",
        "/business/settings",
      ],
    },
    sitemap: new URL("/sitemap.xml", base).toString(),
  };
}
