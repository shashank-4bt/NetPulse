export const ISOLATION_LAYERS = [
  {
    id: "device",
    name: "Device",
    summary: "Local resolver, clock, and client configuration.",
  },
  {
    id: "wifi",
    name: "Wi-Fi",
    summary: "Access point association, signal, and local congestion.",
  },
  {
    id: "dns",
    name: "DNS",
    summary: "Name resolution path, resolver, and NXDOMAIN vs timeout.",
  },
  {
    id: "connectivity",
    name: "Connectivity",
    summary: "Whether packets leave the network at all.",
  },
  {
    id: "isp",
    name: "ISP",
    summary: "Access network and autonomous system conditions.",
  },
  {
    id: "routing",
    name: "Routing",
    summary: "Path changes, loss, and reachability between networks.",
  },
  {
    id: "cdn",
    name: "CDN",
    summary: "Edge POP selection and cache/origin behavior.",
  },
  {
    id: "tls",
    name: "TLS",
    summary: "Certificate, handshake, and protocol negotiation.",
  },
  {
    id: "http",
    name: "HTTP",
    summary: "Status codes, redirects, and application latency.",
  },
  {
    id: "service",
    name: "Service",
    summary: "The destination application itself.",
  },
] as const;
