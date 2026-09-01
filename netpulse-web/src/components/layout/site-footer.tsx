import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { FOOTER_COLUMNS } from "@/lib/content/navigation";
import { publicConfig } from "@/lib/config/public";

export function SiteFooter() {
  return (
    <footer className="mt-auto border-t border-border bg-muted/30">
      <PageContainer className="grid gap-8 py-10 sm:grid-cols-2 lg:grid-cols-4">
        <div>
          <p className="text-sm font-semibold text-foreground">
            {publicConfig.appName}
          </p>
          <p className="mt-2 max-w-xs text-sm text-muted-foreground">
            {publicConfig.appTagline}. Isolates whether a failure is you, your
            network, your ISP, or the service.
          </p>
        </div>
        {FOOTER_COLUMNS.map((column) => (
          <div key={column.title}>
            <p className="text-sm font-medium text-foreground">{column.title}</p>
            <ul className="mt-3 space-y-2">
              {column.links.map((link) => (
                <li key={link.href}>
                  <Link
                    href={link.href}
                    className="text-sm text-muted-foreground hover:text-foreground"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </PageContainer>
      <PageContainer className="border-t border-border py-4 text-xs text-muted-foreground">
        <p>
          Public pages may describe product method. They do not invent live
          measurements, incidents, or service status.
        </p>
      </PageContainer>
    </footer>
  );
}
