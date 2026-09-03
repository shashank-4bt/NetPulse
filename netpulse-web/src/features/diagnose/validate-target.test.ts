import { describe, expect, it } from "vitest";

import { validateDiagnoseTarget } from "@/features/diagnose/validate-target";

describe("validateDiagnoseTarget", () => {
  it("accepts a public hostname", () => {
    expect(validateDiagnoseTarget("example.com")).toMatchObject({
      ok: true,
      hostname: "example.com",
      kind: "domain",
    });
  });

  it("accepts an https URL and extracts the host", () => {
    expect(validateDiagnoseTarget("https://Status.Example.com/path")).toMatchObject({
      ok: true,
      hostname: "status.example.com",
      kind: "url",
    });
  });

  it("accepts known services by name or host", () => {
    expect(validateDiagnoseTarget("youtube.com")).toMatchObject({
      ok: true,
      hostname: "youtube.com",
      kind: "known_service",
      serviceSlug: "youtube",
    });
    expect(validateDiagnoseTarget("Google")).toMatchObject({
      ok: true,
      hostname: "google.com",
      kind: "known_service",
    });
    expect(validateDiagnoseTarget("instagram.com")).toMatchObject({
      ok: true,
      hostname: "instagram.com",
      kind: "known_service",
    });
  });

  it("rejects empty input", () => {
    expect(validateDiagnoseTarget("  ").ok).toBe(false);
  });

  it("rejects unsafe protocols, credentials, and injection characters", () => {
    expect(validateDiagnoseTarget("javascript:alert(1)").ok).toBe(false);
    expect(validateDiagnoseTarget("file:///etc/passwd").ok).toBe(false);
    expect(validateDiagnoseTarget("https://user:pass@example.com").ok).toBe(
      false
    );
    expect(validateDiagnoseTarget("<script>example.com").ok).toBe(false);
  });

  it("rejects localhost and private addresses", () => {
    expect(validateDiagnoseTarget("localhost").ok).toBe(false);
    expect(validateDiagnoseTarget("http://localhost/").ok).toBe(false);
    expect(validateDiagnoseTarget("127.0.0.1").ok).toBe(false);
    expect(validateDiagnoseTarget("10.0.0.1").ok).toBe(false);
    expect(validateDiagnoseTarget("192.168.1.1").ok).toBe(false);
    expect(validateDiagnoseTarget("169.254.169.254").ok).toBe(false);
    expect(validateDiagnoseTarget("http://169.254.169.254/latest/meta-data").ok).toBe(
      false
    );
    expect(validateDiagnoseTarget("::1").ok).toBe(false);
    expect(validateDiagnoseTarget("168.63.129.16").ok).toBe(false);
    expect(validateDiagnoseTarget("100.64.0.1").ok).toBe(false);
    expect(validateDiagnoseTarget("foo.localhost").ok).toBe(false);
    expect(validateDiagnoseTarget("db.internal").ok).toBe(false);
  });
});
