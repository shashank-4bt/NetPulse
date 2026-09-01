import { describe, expect, it } from "vitest";

import { getServiceBySlug, getServiceSlugs, SERVICE_CATALOG } from "@/lib/content/services";

describe("service catalog", () => {
  it("uses unique slugs and can look up every entry", () => {
    const slugs = getServiceSlugs();
    expect(new Set(slugs).size).toBe(SERVICE_CATALOG.length);
    for (const slug of slugs) {
      expect(getServiceBySlug(slug)?.slug).toBe(slug);
    }
  });

  it("does not attach a live status field to catalog entries", () => {
    for (const service of SERVICE_CATALOG) {
      expect(service).not.toHaveProperty("status");
      expect(service).not.toHaveProperty("uptime");
    }
  });
});
