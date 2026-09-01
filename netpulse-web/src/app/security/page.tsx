import type { Metadata } from "next";
import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { PageHero } from "@/components/public/page-hero";
import { Button } from "@/components/ui/button";

export const metadata: Metadata = {
  title: "Security",
  description:
    "How NetPulse treats measurement safety, including SSRF controls and session rules.",
  alternates: { canonical: "/security" },
};

export default function SecurityPage() {
  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Security"
        title="Measurement must not become an open proxy"
        description="Anyone can type a target. Workers must refuse private networks, loopback, link-local addresses, and cloud metadata."
      />
      <PageContainer className="max-w-2xl space-y-8 py-10 text-sm text-muted-foreground md:text-base">
        <section aria-labelledby="ssrf-heading" className="space-y-2">
          <h2 id="ssrf-heading" className="text-lg font-semibold text-foreground">
            SSRF policy
          </h2>
          <p>
            The diagnose form already rejects localhost and private IPv4
            literals. Server-side workers will resolve DNS, check the resulting
            addresses, and re-validate after redirects before connecting.
          </p>
        </section>
        <section aria-labelledby="session-heading" className="space-y-2">
          <h2 id="session-heading" className="text-lg font-semibold text-foreground">
            Sessions
          </h2>
          <p>
            When accounts exist, tokens belong in HTTP-only cookies — not
            localStorage. This public site has no login.
          </p>
        </section>
        <section aria-labelledby="abuse-heading" className="space-y-2">
          <h2 id="abuse-heading" className="text-lg font-semibold text-foreground">
            Abuse
          </h2>
          <p>
            Diagnose endpoints will be rate-limited. The current form does not
            send probes, so it cannot be used to scan third parties from
            NetPulse infrastructure.
          </p>
        </section>
        <Button
          variant="outline"
          nativeButton={false}
          render={<Link href="/developers" />}
        >
          Developers
        </Button>
      </PageContainer>
    </main>
  );
}
