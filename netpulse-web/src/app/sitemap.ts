import type { MetadataRoute } from "next";

import { getServiceSlugs } from "@/lib/content/services";

const STATIC_PATHS = [
  "/",
  "/about",
  "/how-it-works",
  "/services",
  "/status",
  "/outages",
  "/map",
  "/developers",
  "/trust",
  "/privacy",
  "/security",
] as const;

export default function sitemap(): MetadataRoute.Sitemap {
  const base = process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000";
  const now = new Date();

  const staticEntries = STATIC_PATHS.map((path) => ({
    url: new URL(path, base).toString(),
    lastModified: now,
    changeFrequency: "weekly" as const,
  }));

  const serviceEntries = getServiceSlugs().map((slug) => ({
    url: new URL(`/service/${slug}`, base).toString(),
    lastModified: now,
    changeFrequency: "monthly" as const,
  }));

  return [...staticEntries, ...serviceEntries];
}
