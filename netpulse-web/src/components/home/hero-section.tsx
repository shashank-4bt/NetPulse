import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { Button } from "@/components/ui/button";
import { publicConfig } from "@/lib/config/public";

export function HeroSection() {
  return (
    <section
      aria-labelledby="hero-heading"
      className="border-b border-border bg-muted/20"
    >
      <PageContainer className="py-12 md:py-16">
        <p className="text-xs font-medium tracking-wide text-muted-foreground uppercase">
          {publicConfig.appTagline}
        </p>
        <h1
          id="hero-heading"
          className="mt-3 max-w-3xl text-4xl font-semibold tracking-tight md:text-5xl"
        >
          Is the internet down?
        </h1>
        <p className="mt-3 max-w-2xl text-lg text-muted-foreground">
          Find out what&apos;s actually broken.
        </p>
        <p className="mt-4 max-w-2xl text-sm text-muted-foreground md:text-base">
          NetPulse measures the path from your device to a service and
          separates local, ISP, routing, and application failures. It does not
          guess, and it will not invent a status when evidence is missing.
        </p>
        <div className="mt-8 flex flex-wrap gap-3">
          <Button nativeButton={false} render={<Link href="/diagnose" />}>
            Check My Internet
          </Button>
          <Button
            variant="outline"
            nativeButton={false}
            render={<Link href="/status" />}
          >
            View Global Health
          </Button>
        </div>
      </PageContainer>
    </section>
  );
}
