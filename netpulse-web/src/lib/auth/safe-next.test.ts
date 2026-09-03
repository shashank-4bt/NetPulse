import { describe, expect, it } from "vitest";

import { safeNext } from "@/lib/auth/safe-next";

describe("safeNext", () => {
  it("allows internal account paths", () => {
    expect(safeNext("/dashboard")).toBe("/dashboard");
    expect(safeNext("/account/privacy")).toBe("/account/privacy");
    expect(safeNext("/admin/users")).toBe("/admin/users");
  });

  it("rejects protocol-relative and off-site tricks", () => {
    expect(safeNext("//evil.example")).toBe("/dashboard");
    expect(safeNext("/\\evil.example")).toBe("/dashboard");
    expect(safeNext("https://evil.example")).toBe("/dashboard");
    expect(safeNext("/https://evil.example")).toBe("/dashboard");
    expect(safeNext("/login@evil.example")).toBe("/dashboard");
    expect(safeNext("evil.example")).toBe("/dashboard");
  });

  it("rejects paths outside the allowlist", () => {
    expect(safeNext("/login")).toBe("/dashboard");
    expect(safeNext("/")).toBe("/dashboard");
  });
});
