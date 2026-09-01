import type { Metadata } from "next";
import Link from "next/link";

import { EmptyState } from "@/components/feedback/empty-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { getPublicHealthSnapshot } from "@/lib/api/public-health";
import { SERVICE_CATALOG } from "@/lib/content/services";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Status",
  description:
    "Public service status. Live health is unavailable until measurement workers are connected.",
  alternates: { canonical: "/status" },
};

export default async function StatusPage() {
  const health = await getPublicHealthSnapshot();

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Public status"
        title="Service health"
        description="This board lists catalog targets only. It is not a live status page."
      />
      <PageContainer className="space-y-6 py-10">
        <DevelopmentBanner description={health.reason} />
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Service</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Result</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {SERVICE_CATALOG.map((service) => (
              <TableRow key={service.slug}>
                <TableCell>
                  <Link
                    href={`/service/${service.slug}`}
                    className="hover:underline"
                  >
                    {service.name}
                  </Link>
                </TableCell>
                <TableCell>{service.category}</TableCell>
                <TableCell>
                  <Badge variant="outline">Not measured</Badge>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        {health.incidents.length === 0 ? (
          <EmptyState
            title="No live results"
            description="When probes exist, this table will show measured layer outcomes instead of 'Not measured'."
          />
        ) : null}
      </PageContainer>
    </main>
  );
}
