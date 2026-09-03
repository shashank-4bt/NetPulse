import { describe, expect, it, vi } from "vitest";

import { trustedClientIP } from "@/lib/auth/bff";

describe("trustedClientIP", () => {
  it("ignores client-supplied forwarding headers by default", () => {
    const request = new Request("http://localhost/api/auth/login", {
      headers: {
        "x-forwarded-for": "8.8.8.8",
        "x-real-ip": "8.8.8.8",
      },
    });
    expect(trustedClientIP(request)).toBeNull();
  });

  it("reads platform headers only when the web trust-proxy flag is set", () => {
    vi.stubEnv("NETPULSE_WEB_TRUST_PROXY", "true");
    const request = new Request("http://localhost/api/auth/login", {
      headers: { "x-real-ip": "203.0.113.9" },
    });
    expect(trustedClientIP(request)).toBe("203.0.113.9");
    vi.unstubAllEnvs();
  });

  it("rejects forwarding headers that are not a single IP", () => {
    vi.stubEnv("NETPULSE_WEB_TRUST_PROXY", "true");
    const request = new Request("http://localhost/api/auth/login", {
      headers: { "x-forwarded-for": "not-an-ip" },
    });
    expect(trustedClientIP(request)).toBeNull();
    vi.unstubAllEnvs();
  });
});
