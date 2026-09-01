import { describe, expect, it } from "vitest";

import { safeDiagnosePrefill } from "@/features/diagnose/safe-prefill";

describe("safeDiagnosePrefill", () => {
  it("accepts a hostname", () => {
    expect(safeDiagnosePrefill(" youtube.com ")).toBe("youtube.com");
  });

  it("drops injection characters and oversized input", () => {
    expect(safeDiagnosePrefill("<script>")).toBe("");
    expect(safeDiagnosePrefill("a".repeat(2049))).toBe("");
  });
});
