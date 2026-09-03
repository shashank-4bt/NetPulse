import type { ReactNode } from "react";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";

export function AdminUnavailable({ title, description }: { title: string; description: string }) {
  return (
    <main id="main-content" className="flex-1">
      <PageHero eyebrow="Operations" title={title} description={description} />
      <PageContainer className="py-10">
        <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set. No demo operators, health scores, or invented error rates are shown." />
      </PageContainer>
    </main>
  );
}

export function AdminForbidden() {
  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Operations"
        title="This surface is not available."
        description="The requested page could not be found."
      />
    </main>
  );
}

export function AdminScreen({
  title,
  description,
  children,
}: {
  title: string;
  description: string;
  children: ReactNode;
}) {
  return (
    <main id="main-content" className="flex-1">
      <PageHero eyebrow="Operations" title={title} description={description} />
      <PageContainer className="space-y-6 py-10">{children}</PageContainer>
    </main>
  );
}
