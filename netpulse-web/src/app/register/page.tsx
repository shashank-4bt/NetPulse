import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { PageHero } from "@/components/public/page-hero";
import { AuthForm } from "@/features/account/auth-form";
import { redirectIfSignedIn } from "@/lib/account/guard";
import { isApiConfigured } from "@/lib/api/backend";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Create account",
  description: "Register with email and password. Email delivery is not configured.",
  robots: { index: false, follow: false },
};

export default async function RegisterPage() {
  await redirectIfSignedIn();
  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Account"
        title="Create an account"
        description="Email verification is prepared, but a message is not sent until a mailer is configured."
      />
      <PageContainer className="max-w-md space-y-6 py-10">
        {isApiConfigured() ? (
          <AuthForm mode="register" />
        ) : (
          <p className="text-sm text-muted-foreground">
            NETPULSE_API_BASE_URL is not set. Registration is unavailable.
          </p>
        )}
      </PageContainer>
    </main>
  );
}
