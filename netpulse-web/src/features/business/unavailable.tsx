import type { ReactNode } from "react";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CreateOrgForm } from "@/features/business/actions";

export function BusinessUnavailable({
  title,
  description,
}: {
  title: string;
  description: string;
}) {
  return (
    <main id="main-content" className="flex-1">
      <PageHero eyebrow="Business" title={title} description={description} />
      <PageContainer className="py-10">
        <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set. No demo members, devices, or invented availability percentages are shown." />
      </PageContainer>
    </main>
  );
}

export function BusinessScreen({
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
      <PageHero eyebrow="Business" title={title} description={description} />
      <PageContainer className="space-y-6 py-10">{children}</PageContainer>
    </main>
  );
}

export function NoOrganization() {
  return (
    <BusinessScreen
      title="Create an organization"
      description="Membership is required before any organization resource can be read."
    >
      <Card>
        <CardHeader>
          <CardTitle>New organization</CardTitle>
        </CardHeader>
        <CardContent>
          <CreateOrgForm />
        </CardContent>
      </Card>
    </BusinessScreen>
  );
}
