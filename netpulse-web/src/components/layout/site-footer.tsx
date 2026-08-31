import { PageContainer } from "@/components/layout/page-container";
import { publicConfig } from "@/lib/config/public";

export function SiteFooter() {
  return (
    <footer className="mt-auto border-t border-border bg-muted/30">
      <PageContainer className="flex flex-col gap-2 py-6 text-sm text-muted-foreground sm:flex-row sm:items-center sm:justify-between">
        <p>
          {publicConfig.appName} — {publicConfig.appTagline}
        </p>
        <p>Design system foundation. Product diagnosis is not available yet.</p>
      </PageContainer>
    </footer>
  );
}
