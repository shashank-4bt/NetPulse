import { afterEach, describe, expect, it, vi } from "vitest";

import { POST } from "@/app/api/auth/[...path]/route";

describe("auth BFF", () => {
  afterEach(() => {
    vi.unstubAllEnvs();
  });

  it("does not invent an account when the API is unset", async () => {
    vi.stubEnv("NETPULSE_API_BASE_URL", "");
    const response = await POST(
      new Request("http://localhost:3000/api/auth/login", {
        method: "POST",
        headers: {
          origin: "http://localhost:3000",
          "content-type": "application/json",
        },
        body: JSON.stringify({ email: "a@example.com", password: "correct-horse" }),
      }),
      { params: Promise.resolve({ path: ["login"] }) }
    );
    expect(response.status).toBe(503);
    const body = (await response.json()) as { ok: boolean; error?: { code?: string } };
    expect(body.ok).toBe(false);
    expect(body.error?.code).toBe("unavailable");
  });
});
