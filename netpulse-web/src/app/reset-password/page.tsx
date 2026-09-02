import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { PageHero } from "@/components/public/page-hero";
import { AuthForm } from "@/features/account/auth-form";
import { isApiConfigured } from "@/lib/api/backend";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Set a new password",
  robots: { index: false, follow: false },
};

type PageProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export default async function ResetPasswordPage({ searchParams }: PageProps) {
  const params = await searchParams;
  const token = typeof params.token === "string" ? params.token : "";
  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Account"
        title="Set a new password"
        description="Use the reset token issued for your account. Other sessions are revoked after a successful reset."
      />
      <PageContainer className="max-w-md space-y-6 py-10">
        {isApiConfigured() ? (
          <AuthForm mode="reset" token={token} />
        ) : (
          <p className="text-sm text-muted-foreground">
            NETPULSE_API_BASE_URL is not set. Password reset is unavailable.
          </p>
        )}
      </PageContainer>
    </main>
  );
}
