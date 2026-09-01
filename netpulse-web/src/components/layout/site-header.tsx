import Link from "next/link";

import { MobileNav } from "@/components/layout/mobile-nav";
import { PageContainer } from "@/components/layout/page-container";
import { ThemeToggle } from "@/components/layout/theme-toggle";
import { Button } from "@/components/ui/button";
import { MOBILE_NAV, PRIMARY_NAV } from "@/lib/content/navigation";
import { publicConfig } from "@/lib/config/public";

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-40 border-b border-border bg-background/95 backdrop-blur-sm">
      <PageContainer className="flex h-14 items-center justify-between gap-3">
        <Link
          href="/"
          className="shrink-0 text-sm font-semibold tracking-tight text-foreground"
        >
          {publicConfig.appName}
        </Link>
        <nav aria-label="Primary" className="hidden items-center gap-1 lg:flex">
          {PRIMARY_NAV.map((item) => (
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
          <Button
            nativeButton={false}
            className="hidden sm:inline-flex"
            render={<Link href="/#diagnose" />}
          >
            Check My Internet
          </Button>
          <ThemeToggle />
          <MobileNav links={MOBILE_NAV} />
        </div>
      </PageContainer>
    </header>
  );
}
