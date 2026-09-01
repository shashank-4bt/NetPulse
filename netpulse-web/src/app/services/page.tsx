import type { Metadata } from "next";
import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { SERVICE_CATALOG } from "@/lib/content/services";

export const metadata: Metadata = {
  title: "Services",
  description:
    "Catalog of services NetPulse can diagnose. This is not a live status board.",
  alternates: { canonical: "/services" },
};

export default function ServicesPage() {
  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Services"
        title="What can be diagnosed"
        description="These entries describe targets and layers. They do not report current availability."
      />
      <PageContainer className="space-y-6 py-10">
        <DevelopmentBanner description="Service pages are editorial. Live health requires measurement workers." />
        <ul className="grid gap-3 sm:grid-cols-2">
          {SERVICE_CATALOG.map((service) => (
            <li key={service.slug}>
              <Card className="h-full">
                <CardHeader>
                  <CardTitle>
                    <Link
                      href={`/service/${service.slug}`}
                      className="hover:underline"
                    >
                      {service.name}
                    </Link>
                  </CardTitle>
                  <CardDescription>{service.category}</CardDescription>
                </CardHeader>
                <CardContent className="space-y-3">
                  <p className="text-sm text-muted-foreground">
                    {service.summary}
                  </p>
                  <Badge variant="outline">Not measured</Badge>
                </CardContent>
              </Card>
            </li>
          ))}
        </ul>
      </PageContainer>
    </main>
  );
}
