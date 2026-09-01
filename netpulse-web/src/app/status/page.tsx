import type { Metadata } from "next";
import Link from "next/link";

import { EmptyState } from "@/components/feedback/empty-state";
import { ErrorState } from "@/components/feedback/error-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { LayerStatusBadge } from "@/components/status/layer-status-badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { emptyServiceIntelligence } from "@/features/observatory/empty-intelligence";
import { loadServiceCatalog } from "@/lib/observatory/load";

export const dynamic = "force-dynamic";

export async function generateMetadata(): Promise<Metadata> {
  const catalog = await loadServiceCatalog();
  return {
    title: "Status",
    description:
      catalog.state === "unavailable"
        ? "Public service status. Live health is unavailable until measurement workers are connected."
        : "Public service status from stored measurements only. Unmeasured rows stay not measured.",
    alternates: { canonical: "/status" },
  };
}

export default async function StatusPage() {
  const catalog = await loadServiceCatalog();
  const intelligence = emptyServiceIntelligence();

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Public status"
        title="Service health"
        description="This board lists catalog targets. It is not a live status page and does not invent operational green."
      />
      <PageContainer className="space-y-6 py-10">
        <DevelopmentBanner description={catalog.reason} />
        {catalog.state === "error" ? (
          <ErrorState title="Status source failed" description={catalog.reason} />
        ) : null}
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Service</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Current state</TableHead>
              <TableHead>Health</TableHead>
              <TableHead>Last updated</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {catalog.services.map((service) => (
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
                  <LayerStatusBadge status={intelligence.currentState === "not_measured" ? "not_measured" : "insufficient_evidence"} />
                </TableCell>
                <TableCell className="text-muted-foreground">
                  Not scored
                </TableCell>
                <TableCell className="font-mono text-muted-foreground">
                  —
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <EmptyState
          title="No live results"
          description="When probe series exist, this table will show measured layer outcomes instead of 'Not measured'."
        />
      </PageContainer>
    </main>
  );
}
