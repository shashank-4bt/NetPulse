export const PRIMARY_NAV = [
  { href: "/diagnose", label: "Diagnose" },
  { href: "/how-it-works", label: "How it works" },
  { href: "/services", label: "Services" },
  { href: "/status", label: "Status" },
  { href: "/developers", label: "Developers" },
  { href: "/trust", label: "Trust" },
] as const;

export const MOBILE_NAV = [
  { href: "/", label: "Home" },
  ...PRIMARY_NAV,
  { href: "/about", label: "About" },
  { href: "/outages", label: "Outages" },
  { href: "/map", label: "Map" },
  { href: "/privacy", label: "Privacy" },
  { href: "/security", label: "Security" },
] as const;

export const FOOTER_COLUMNS = [
  {
    title: "Product",
    links: [
      { href: "/diagnose", label: "Diagnose" },
      { href: "/how-it-works", label: "How it works" },
      { href: "/services", label: "Services" },
      { href: "/status", label: "Status" },
      { href: "/outages", label: "Outages" },
      { href: "/map", label: "Map" },
    ],
  },
  {
    title: "Company",
    links: [
      { href: "/about", label: "About" },
      { href: "/developers", label: "Developers" },
      { href: "/trust", label: "Trust" },
    ],
  },
  {
    title: "Legal",
    links: [
      { href: "/privacy", label: "Privacy" },
      { href: "/security", label: "Security" },
    ],
  },
] as const;
