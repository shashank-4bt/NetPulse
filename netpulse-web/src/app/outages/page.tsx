import type { Metadata } from "next";

import { EmptyState } from "@/components/feedback/empty-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { getPublicHealthSnapshot } from "@/lib/api/public-health";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Outages",
  description:
    "Confirmed internet incidents. The feed is not connected; no outages are invented.",
  alternates: { canonical: "/outages" },
};

export default async function OutagesPage() {
  const health = await getPublicHealthSnapshot();

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Outages"
        title="Incident feed"
        description="An empty feed is not a claim that the internet is healthy. NetPulse publishes incidents only after evidence exists."
      />
      <PageContainer className="space-y-6 py-10">
        <DevelopmentBanner description={health.reason} />
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Time</TableHead>
              <TableHead>Scope</TableHead>
              <TableHead>Evidence</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell colSpan={3} className="text-muted-foreground">
                No incident records.
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
        <EmptyState
          title="No outages on record"
          description="When the store is connected, each row will cite measurements and a confidence level."
        />
      </PageContainer>
    </main>
  );
}
