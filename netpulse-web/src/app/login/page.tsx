import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { PageHero } from "@/components/public/page-hero";
import { AuthForm } from "@/features/account/auth-form";
import { redirectIfSignedIn } from "@/lib/account/guard";
import { isApiConfigured } from "@/lib/api/backend";
import { safeNext } from "@/lib/auth/safe-next";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Sign in",
  description: "Sign in to NetPulse. Sessions use an HTTP-only cookie.",
  robots: { index: false, follow: false },
};

type PageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export default async function LoginPage({ searchParams }: PageProps) {
  await redirectIfSignedIn();
  const params = await searchParams;
  const next = typeof params.next === "string" ? params.next : "/dashboard";

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Account"
        title="Sign in"
        description="Password sign-in only. OAuth, passkeys, and MFA are not configured."
      />
      <PageContainer className="max-w-md space-y-6 py-10">
        {isApiConfigured() ? (
          <AuthForm mode="login" nextPath={safeNext(next)} />
        ) : (
          <p className="text-sm text-muted-foreground">
            NETPULSE_API_BASE_URL is not set. Accounts are unavailable. No local
            demo login is offered.
          </p>
        )}
      </PageContainer>
    </main>
  );
}
