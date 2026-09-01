import { describe, expect, it } from "vitest";

import { validateDiagnoseTarget } from "@/features/diagnose/validate-target";

describe("validateDiagnoseTarget", () => {
  it("accepts a public hostname", () => {
    expect(validateDiagnoseTarget("example.com")).toEqual({
      ok: true,
      hostname: "example.com",
    });
  });

  it("accepts an https URL and extracts the host", () => {
    expect(validateDiagnoseTarget("https://Status.Example.com/path")).toEqual({
      ok: true,
      hostname: "status.example.com",
    });
  });

  it("rejects empty input", () => {
    expect(validateDiagnoseTarget("  ").ok).toBe(false);
  });

  it("rejects localhost and private addresses", () => {
    expect(validateDiagnoseTarget("localhost").ok).toBe(false);
    expect(validateDiagnoseTarget("127.0.0.1").ok).toBe(false);
    expect(validateDiagnoseTarget("10.0.0.1").ok).toBe(false);
    expect(validateDiagnoseTarget("192.168.1.1").ok).toBe(false);
    expect(validateDiagnoseTarget("169.254.169.254").ok).toBe(false);
  });
});
