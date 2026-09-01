import type { Metadata } from "next";

import { ChartContainer } from "@/components/data/chart-container";
import { EvidenceItem } from "@/components/data/evidence-item";
import { MetricCard } from "@/components/data/metric-card";
import { Timeline } from "@/components/data/timeline";
import { InteractiveControls } from "@/components/design-system/interactive-controls";
import { EmptyState } from "@/components/feedback/empty-state";
import { ErrorState } from "@/components/feedback/error-state";
import { InsufficientEvidenceState } from "@/components/feedback/insufficient-evidence-state";
import { LoadingState } from "@/components/feedback/loading-state";
import { UnavailableState } from "@/components/feedback/unavailable-state";
import { PageContainer } from "@/components/layout/page-container";
import { SectionContainer } from "@/components/layout/section-container";
import { ConfidenceBadge } from "@/components/status/confidence-badge";
import { HealthScore } from "@/components/status/health-score";
import { SeverityBadge } from "@/components/status/severity-badge";
import { StatusBadge } from "@/components/status/status-badge";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  CONFIDENCE_LEVELS,
  OPERATIONAL_STATUSES,
  SEVERITIES,
} from "@/domain/display";

export const metadata: Metadata = {
  title: "Components",
  description: "NetPulse reusable component gallery. No live product data.",
  robots: { index: false, follow: false },
};

export default function DesignSystemPage() {
  return (
    <main id="main-content" className="flex-1">
      <PageContainer>
        <SectionContainer labelledBy="ds-heading">
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbLink href="/">Foundations</BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbPage>Components</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
          <h1 id="ds-heading" className="mt-4 text-3xl font-semibold tracking-tight">
            Component gallery
          </h1>
          <p className="mt-2 max-w-2xl text-muted-foreground">
            Interactive examples for the design system. Values here are
            fixtures for UI review. They are not measurements, incidents, or
            live service status.
          </p>
        </SectionContainer>

        <SectionContainer labelledBy="actions-heading">
          <h2 id="actions-heading" className="text-lg font-medium">
            Actions and forms
          </h2>
          <div className="mt-4 flex flex-wrap gap-2">
            <Button type="button">Primary</Button>
            <Button type="button" variant="secondary">
              Secondary
            </Button>
            <Button type="button" variant="outline">
              Outline
            </Button>
            <Button type="button" variant="ghost">
              Ghost
            </Button>
            <Button type="button" variant="destructive">
              Destructive
            </Button>
            <Button type="button" disabled>
              Disabled
            </Button>
          </div>
          <div className="mt-6">
            <InteractiveControls />
          </div>
        </SectionContainer>

        <SectionContainer labelledBy="status-heading">
          <h2 id="status-heading" className="text-lg font-medium">
            Status, severity, confidence
          </h2>
          <p className="mt-1 text-sm text-muted-foreground">
            Each indicator includes an icon and a text label. Color is never
            the only signal.
          </p>
          <div className="mt-4 flex flex-wrap gap-2">
            {OPERATIONAL_STATUSES.map((status) => (
              <StatusBadge key={status} status={status} />
            ))}
          </div>
          <div className="mt-3 flex flex-wrap gap-2">
            {SEVERITIES.map((severity) => (
              <SeverityBadge key={severity} severity={severity} />
            ))}
          </div>
          <div className="mt-3 flex flex-wrap gap-2">
            {CONFIDENCE_LEVELS.map((confidence) => (
              <ConfidenceBadge key={confidence} confidence={confidence} />
            ))}
          </div>
          <div className="mt-6 grid gap-4 md:grid-cols-2">
            <HealthScore value={null} confidence={null} />
            <HealthScore
              value={72}
              confidence="medium"
              label="Health score (fixture)"
            />
          </div>
        </SectionContainer>

        <SectionContainer labelledBy="feedback-heading">
          <h2 id="feedback-heading" className="text-lg font-medium">
            Feedback states
          </h2>
          <div className="mt-4 grid gap-4 lg:grid-cols-2">
            <EmptyState
              title="No diagnostic runs"
              description="Runs appear here after a measurement is requested. None exist yet."
            />
            <ErrorState
              title="Request failed"
              description="The API returned an error. No fallback health data was invented."
            />
            <LoadingState />
            <UnavailableState />
            <InsufficientEvidenceState />
            <Alert>
              <AlertTitle>Alert</AlertTitle>
              <AlertDescription>
                Use alerts for recoverable notices. Do not use them to imply
                an outage without evidence.
              </AlertDescription>
            </Alert>
          </div>
        </SectionContainer>

        <SectionContainer labelledBy="data-heading">
          <h2 id="data-heading" className="text-lg font-medium">
            Data display
          </h2>
          <div className="mt-4 grid gap-4 md:grid-cols-2">
            <MetricCard
              title="Latency"
              description="Example with no measurement"
              value={null}
            />
            <MetricCard
              title="Packet loss"
              description="Fixture for layout only"
              value="—"
              caption="Em dash indicates the gallery, not a measured zero."
            />
          </div>
          <div className="mt-4 grid gap-4 lg:grid-cols-3">
            <EvidenceItem
              evidenceClass="measured_fact"
              title="TLS handshake completed"
              body="A measured observation belongs here only when a probe recorded it."
              source="probe:tls (fixture copy)"
            />
            <EvidenceItem
              evidenceClass="inferred_hypothesis"
              title="ISP congestion possible"
              body="Observed evidence suggests a hypothesis. This is not a measured fact."
            />
            <EvidenceItem
              evidenceClass="recommendation"
              title="Re-test after DNS change"
              body="Recommendations are actions. They are not evidence of a failure."
            />
          </div>
          <div className="mt-4">
            <Timeline
              events={[
                {
                  id: "1",
                  timestampLabel: "Fixture",
                  title: "Job queued",
                  detail: "Timeline copy for layout. Not a real run.",
                },
                {
                  id: "2",
                  timestampLabel: "Fixture",
                  title: "Worker unavailable",
                  detail: "Measurement workers are not connected in Stage 02.",
                },
              ]}
            />
          </div>
          <div className="mt-4">
            <ChartContainer
              title="Latency series"
              description="Charts render only when a measured series is provided."
              state="empty"
            />
          </div>
        </SectionContainer>

        <SectionContainer labelledBy="chrome-heading">
          <h2 id="chrome-heading" className="text-lg font-medium">
            Structure
          </h2>
          <Tabs defaultValue="overview" className="mt-4">
            <TabsList>
              <TabsTrigger value="overview">Overview</TabsTrigger>
              <TabsTrigger value="evidence">Evidence</TabsTrigger>
            </TabsList>
            <TabsContent value="overview" className="mt-3">
              Tabs keep related panels on one page without fake metrics.
            </TabsContent>
            <TabsContent value="evidence" className="mt-3">
              Evidence panels must preserve fact vs hypothesis labels.
            </TabsContent>
          </Tabs>
          <Card className="mt-4">
            <CardHeader>
              <CardTitle>Card</CardTitle>
              <CardDescription>Surface for grouped content.</CardDescription>
            </CardHeader>
            <CardContent className="flex flex-wrap gap-2">
              <Badge>Default</Badge>
              <Badge variant="secondary">Secondary</Badge>
              <Badge variant="outline">Outline</Badge>
              <Skeleton className="h-4 w-24" />
            </CardContent>
          </Card>
          <div className="mt-4 overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Layer</TableHead>
                  <TableHead>Result</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow>
                  <TableCell>DNS</TableCell>
                  <TableCell>Not measured</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell>TLS</TableCell>
                  <TableCell>Not measured</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </div>
          <Pagination className="mt-4">
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious href="/design-system" />
              </PaginationItem>
              <PaginationItem>
                <PaginationNext href="/design-system" />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </SectionContainer>
      </PageContainer>
    </main>
  );
}
