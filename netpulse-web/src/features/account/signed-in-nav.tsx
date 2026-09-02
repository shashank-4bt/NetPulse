import Link from "next/link";

import { ACCOUNT_NAV, DASHBOARD_NAV } from "@/domain/account";
import { BUSINESS_NAV } from "@/domain/business";
import { DEVELOPER_NAV } from "@/domain/developer";

type SignedInNavProps = {
  kind: "dashboard" | "account" | "developer" | "business";
};

export function SignedInNav({ kind }: SignedInNavProps) {
  const items =
    kind === "dashboard"
      ? DASHBOARD_NAV
      : kind === "account"
        ? ACCOUNT_NAV
        : kind === "developer"
          ? DEVELOPER_NAV
          : BUSINESS_NAV;
  const label =
    kind === "dashboard"
      ? "Dashboard"
      : kind === "account"
        ? "Account"
        : kind === "developer"
          ? "Developer workspace"
          : "Business workspace";
  return (
    <nav aria-label={label} className="flex flex-wrap gap-2">
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
