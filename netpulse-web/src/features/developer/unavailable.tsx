import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";

export function DeveloperUnavailable({
  title,
  description,
}: {
  title: string;
  description: string;
}) {
  return (
    <main id="main-content" className="flex-1">
      <PageHero eyebrow="Developers" title={title} description={description} />
      <PageContainer className="py-10">
        <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set. No demo monitors, keys, or SLA percentages are shown." />
      </PageContainer>
    </main>
  );
}
