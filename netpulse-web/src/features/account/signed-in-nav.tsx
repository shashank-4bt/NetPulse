import Link from "next/link";

import { ACCOUNT_NAV, DASHBOARD_NAV } from "@/domain/account";

type SignedInNavProps = {
  kind: "dashboard" | "account";
};

export function SignedInNav({ kind }: SignedInNavProps) {
  const items = kind === "dashboard" ? DASHBOARD_NAV : ACCOUNT_NAV;
  return (
    <nav aria-label={kind === "dashboard" ? "Dashboard" : "Account"} className="flex flex-wrap gap-2">
      {items.map((item) => (
        <Link
          key={item.href}
          href={item.href}
          className="inline-flex min-h-11 items-center rounded-md px-3 text-sm text-muted-foreground hover:bg-muted hover:text-foreground"
        >
          {item.label}
        </Link>
      ))}
    </nav>
  );
}
