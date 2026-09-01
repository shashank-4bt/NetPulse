import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { SectionContainer } from "@/components/layout/section-container";
import { SectionHeading } from "@/components/public/section-heading";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { SERVICE_CATALOG } from "@/lib/content/services";

export function ServicesSection() {
  return (
    <SectionContainer labelledBy="services-heading" className="bg-muted/20">
      <PageContainer>
        <SectionHeading
          id="services-heading"
          eyebrow="Popular services"
          title="What NetPulse can diagnose"
          description="This is a catalog, not a live status board. Each service page explains the path we measure."
        />
        <ul className="mt-6 grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
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
                <CardContent>
                  <p className="text-sm text-muted-foreground">
                    {service.summary}
                  </p>
                  <Badge variant="outline" className="mt-3">
                    Not measured
                  </Badge>
                </CardContent>
              </Card>
            </li>
          ))}
        </ul>
        <Button
          className="mt-6"
          variant="outline"
          nativeButton={false}
          render={<Link href="/services" />}
        >
          All services
        </Button>
      </PageContainer>
    </SectionContainer>
  );
}
