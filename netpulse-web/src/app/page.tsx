import type { Metadata } from "next";
import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { SectionContainer } from "@/components/layout/section-container";
import { Button } from "@/components/ui/button";
import { publicConfig } from "@/lib/config/public";

export const metadata: Metadata = {
  title: "Foundations",
  description:
    "NetPulse visual foundations. This is not a live internet health dashboard.",
};

export default function HomePage() {
  return (
    <main id="main-content" className="flex-1">
      <PageContainer>
        <SectionContainer labelledBy="hero-heading">
          <p className="text-xs font-medium tracking-wide text-muted-foreground uppercase">
            Stage 02 · Design system
          </p>
          <h1
            id="hero-heading"
            className="mt-2 max-w-2xl text-3xl font-semibold tracking-tight text-foreground md:text-4xl"
          >
            {publicConfig.appName}
          </h1>
          <p className="mt-3 max-w-2xl text-base text-muted-foreground">
            {publicConfig.appTagline}. This site currently exposes the visual
            language and reusable components only. Diagnostic measurement,
            outage data, and service status are not implemented and are not
            simulated.
          </p>
          <div className="mt-6 flex flex-wrap gap-3">
            <Button nativeButton={false} render={<Link href="/design-system" />}>
              View components
            </Button>
          </div>
        </SectionContainer>
        <SectionContainer labelledBy="principles-heading">
          <h2 id="principles-heading" className="text-lg font-medium">
            Visual principles
          </h2>
          <ul className="mt-4 grid gap-3 sm:grid-cols-2">
            <li className="rounded-lg border border-border bg-card p-4 text-sm">
              <strong className="block font-medium">Precise</strong>
              <span className="text-muted-foreground">
                Infrastructure typography, tabular numbers, restrained surfaces.
              </span>
            </li>
            <li className="rounded-lg border border-border bg-card p-4 text-sm">
              <strong className="block font-medium">Not color-only</strong>
              <span className="text-muted-foreground">
                Status, severity, and confidence always include an icon and a
                text label.
              </span>
            </li>
            <li className="rounded-lg border border-border bg-card p-4 text-sm">
              <strong className="block font-medium">Honest states</strong>
              <span className="text-muted-foreground">
                Empty, error, loading, unavailable, and insufficient evidence
                are first-class.
              </span>
            </li>
            <li className="rounded-lg border border-border bg-card p-4 text-sm">
              <strong className="block font-medium">Calm</strong>
              <span className="text-muted-foreground">
                No neon gradients, no fake live metrics, no decorative charts.
              </span>
            </li>
          </ul>
        </SectionContainer>
      </PageContainer>
    </main>
  );
}
