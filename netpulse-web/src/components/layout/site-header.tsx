import Link from "next/link";

import { MobileNav } from "@/components/layout/mobile-nav";
import { PageContainer } from "@/components/layout/page-container";
import { ThemeToggle } from "@/components/layout/theme-toggle";
import { publicConfig } from "@/lib/config/public";

const NAV_LINKS = [
  { href: "/", label: "Foundations" },
  { href: "/design-system", label: "Components" },
] as const;

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-40 border-b border-border bg-background/95 backdrop-blur-sm">
      <PageContainer className="flex h-14 items-center justify-between gap-3">
        <div className="flex min-w-0 items-center gap-3">
          <Link
            href="/"
            className="truncate text-sm font-semibold tracking-tight text-foreground"
          >
            {publicConfig.appName}
          </Link>
          <p className="hidden truncate text-xs text-muted-foreground md:block">
            {publicConfig.appTagline}
          </p>
        </div>
        <nav aria-label="Primary" className="hidden items-center gap-1 md:flex">
          {NAV_LINKS.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className="inline-flex min-h-11 items-center rounded-md px-3 text-sm text-muted-foreground hover:bg-muted hover:text-foreground"
            >
              {item.label}
            </Link>
          ))}
        </nav>
        <div className="flex items-center gap-2">
          <ThemeToggle />
          <MobileNav links={NAV_LINKS} />
        </div>
      </PageContainer>
    </header>
  );
}

export { NAV_LINKS };
