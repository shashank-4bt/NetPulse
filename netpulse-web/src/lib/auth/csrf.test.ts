import { describe, expect, it } from "vitest";

import { csrfAllowed } from "@/lib/auth/csrf";

describe("csrfAllowed", () => {
  it("allows GET without Origin", () => {
    expect(csrfAllowed(new Request("http://localhost:3000/api/auth/me"))).toBe(true);
  });

  it("rejects cross-origin POST", () => {
    const request = new Request("http://localhost:3000/api/auth/login", {
      method: "POST",
      headers: { origin: "https://evil.example" },
    });
    expect(csrfAllowed(request)).toBe(false);
  });

  it("allows same-origin POST", () => {
    const request = new Request("http://localhost:3000/api/auth/login", {
      method: "POST",
      headers: { origin: "http://localhost:3000" },
    });
    expect(csrfAllowed(request)).toBe(true);
  });
});
