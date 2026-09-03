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
            Diagnose input and workers deny localhost, loopback, private ranges,
            link-local, CGNAT, documentation prefixes, and cloud metadata
            (including Azure IMDS). DNS is resolved, every address is checked,
            connections pin a public IP, and redirects are revalidated so a
            public hostname cannot bounce into a private one.
          </p>
        </section>
        <section aria-labelledby="session-heading" className="space-y-2">
          <h2 id="session-heading" className="text-lg font-semibold text-foreground">
            Sessions and browser security
          </h2>
          <p>
            Sessions use the HTTP-only <code>np_session</code> cookie. Tokens
            are not stored in localStorage. Cookie-authenticated mutations
            require a same-origin request. Pages send Content-Security-Policy,
            frame denial, and related headers. OAuth, passkeys, and MFA
            endpoints exist but stay unavailable until a provider is configured.
          </p>
        </section>
        <section aria-labelledby="abuse-heading" className="space-y-2">
          <h2 id="abuse-heading" className="text-lg font-semibold text-foreground">
            Abuse
          </h2>
          <p>
            Diagnose and authentication endpoints are rate-limited by the
            connecting address. Client-supplied <code>X-Forwarded-For</code> is
            not trusted unless a reverse proxy is explicitly configured to
            overwrite it.
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
