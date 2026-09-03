import { describe, expect, it } from "vitest";

import { sessionCookieOptions } from "@/lib/auth/session";

describe("sessionCookieOptions", () => {
  it("keeps the session cookie HTTP-only and SameSite=Lax", () => {
    const options = sessionCookieOptions();
    expect(options.httpOnly).toBe(true);
    expect(options.sameSite).toBe("lax");
    expect(options.path).toBe("/");
  });
});
