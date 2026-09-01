import type { Metadata } from "next";
import Link from "next/link";

import { EmptyState } from "@/components/feedback/empty-state";
import { ErrorState } from "@/components/feedback/error-state";
import { UnavailableState } from "@/components/feedback/unavailable-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { SeverityBadge } from "@/components/status/severity-badge";
import { Badge } from "@/components/ui/badge";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { SEVERITIES } from "@/domain/display";
import { INCIDENT_STAGE_LABELS, INCIDENT_STAGES, PAGE_SIZE } from "@/domain/observatory";
import { OutageFilters } from "@/features/observatory/outage-filters";
import { outageQueryString } from "@/features/observatory/filter-incidents";
import { SERVICE_CATALOG } from "@/lib/content/services";
import { loadOutages } from "@/lib/observatory/load";
import type { Severity } from "@/domain/display";

export const dynamic = "force-dynamic";

type OutagesPageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export async function generateMetadata({
  searchParams,
}: OutagesPageProps): Promise<Metadata> {
  const loaded = await loadOutages(await searchParams);
  return {
    title: "Outages",
    description:
      loaded.total === 0
        ? "No stored internet incidents. An empty feed is not a healthy-internet claim."
        : `${loaded.total} stored incident records. Each row is a stored event, not a population impact estimate.`,
    alternates: { canonical: "/outages" },
  };
}

export default async function OutagesPage({ searchParams }: OutagesPageProps) {
  const loaded = await loadOutages(await searchParams);
  const pageCount = Math.max(1, Math.ceil(loaded.total / PAGE_SIZE));
  const query = loaded.query;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Outages"
        title="Incident feed"
        description="An empty feed is not a claim that the internet is healthy. NetPulse publishes incidents only after evidence exists."
      />
      <PageContainer className="space-y-6 py-10">
        <DevelopmentBanner description={loaded.reason} />
        <OutageFilters query={query} services={SERVICE_CATALOG} />

        {loaded.state === "unavailable" ? (
          <UnavailableState
            title="Incident store unavailable"
            description={loaded.reason}
          />
        ) : null}
        {loaded.state === "error" ? (
          <ErrorState title="Incident feed failed" description={loaded.reason} />
        ) : null}

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Started</TableHead>
              <TableHead>Title</TableHead>
              <TableHead>Severity</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Scope</TableHead>
              <TableHead>Sample</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loaded.items.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-muted-foreground">
                  No incident records match these filters.
                </TableCell>
              </TableRow>
            ) : (
              loaded.items.map((incident) => (
                <TableRow key={incident.id}>
                  <TableCell className="font-mono text-xs">
                    {incident.startedAt || "—"}
                  </TableCell>
                  <TableCell>
                    <Link
                      href={`/incident/${incident.id}`}
                      className="hover:underline"
                    >
                      {incident.title}
                    </Link>
                  </TableCell>
                  <TableCell>
                    {isSeverity(incident.severity) ? (
                      <SeverityBadge severity={incident.severity} />
                    ) : (
                      <Badge variant="outline">Unclassified</Badge>
                    )}
                  </TableCell>
                  <TableCell>
                    {isStage(incident.status)
                      ? INCIDENT_STAGE_LABELS[incident.status]
                      : incident.status || "Unclassified"}
                  </TableCell>
                  <TableCell>{incident.scope || "—"}</TableCell>
                  <TableCell className="font-mono text-xs">
                    {incident.sampleCount}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>

        {loaded.total === 0 && loaded.state !== "error" ? (
          <EmptyState
            title="No outages on record"
            description="When the store contains incidents, each row will cite evidence, sample count, and a confidence caveat. Filters, search, sort, and pagination still apply to stored rows only."
          />
        ) : null}

        {loaded.total > PAGE_SIZE ? (
          <Pagination>
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  href={
                    query.page > 1
                      ? `/outages${outageQueryString(query, query.page - 1)}`
                      : undefined
                  }
                  aria-disabled={query.page <= 1}
                />
              </PaginationItem>
              <PaginationItem>
                <PaginationLink href={`/outages${outageQueryString(query, query.page)}`} isActive>
                  {query.page}
                </PaginationLink>
              </PaginationItem>
              <PaginationItem>
                <PaginationNext
                  href={
                    query.page < pageCount
                      ? `/outages${outageQueryString(query, query.page + 1)}`
                      : undefined
                  }
                  aria-disabled={query.page >= pageCount}
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        ) : null}
      </PageContainer>
    </main>
  );
}

function isSeverity(value: string): value is Severity {
  return (SEVERITIES as readonly string[]).includes(value);
}

function isStage(
  value: string
): value is (typeof INCIDENT_STAGES)[number] {
  return (INCIDENT_STAGES as readonly string[]).includes(value);
}
