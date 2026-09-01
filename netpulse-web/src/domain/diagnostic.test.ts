import { describe, expect, it } from "vitest";

import {
  DIAGNOSTIC_OUTCOMES,
  DIAGNOSTIC_STEP_IDS,
  EVIDENCE_GRAPH_NODE_IDS,
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

  it("keeps the Stage 05 evidence graph order", () => {
    expect(EVIDENCE_GRAPH_NODE_IDS).toEqual([
      "device",
      "wifi",
      "router",
      "isp",
      "route",
      "cdn",
      "service",
    ]);
  });

  it("lists technical detail blocks including route metadata and network/ASN", () => {
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
    expect(MEASUREMENT_BLOCKS.map((block) => block.label)).toEqual([
      "DNS",
      "TCP",
      "TLS",
      "HTTP",
      "Latency",
      "Packet loss",
      "Route metadata",
      "Network/ASN",
      "Region",
      "Timestamp",
    ]);
  });
});
