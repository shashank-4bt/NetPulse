import { describe, expect, it } from "vitest";

import {
  DIAGNOSTIC_OUTCOMES,
  DIAGNOSTIC_STEP_IDS,
  MEASUREMENT_BLOCKS,
} from "@/domain/diagnostic";

describe("diagnostic domain", () => {
  it("keeps the Stage 04 step order", () => {
    expect(DIAGNOSTIC_STEP_IDS).toEqual([
      "initializing",
      "device",
      "wifi",
      "dns",
      "connectivity",
      "isp",
      "routing",
      "tls",
      "http",
      "service",
      "regional_comparison",
      "network_comparison",
      "analysis",
      "complete",
    ]);
  });

  it("includes every required outcome without mapping unknown to failure", () => {
    expect(DIAGNOSTIC_OUTCOMES).toEqual([
      "success",
      "partial_success",
      "timeout",
      "backend_unavailable",
      "insufficient_evidence",
      "invalid_input",
      "measurement_unavailable",
    ]);
  });

  it("lists technical detail blocks without charts as a requirement", () => {
    expect(MEASUREMENT_BLOCKS.map((block) => block.key)).toEqual([
      "dns",
      "tcp",
      "tls",
      "http",
      "latency",
      "packet_loss",
      "routing",
      "network",
      "region",
      "timestamp",
    ]);
  });
});
