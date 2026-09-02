import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { PageHero } from "@/components/public/page-hero";
import { AuthForm } from "@/features/account/auth-form";
import { isApiConfigured } from "@/lib/api/backend";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Reset password",
  description: "Request a password reset. The response does not reveal whether an account exists.",
  robots: { index: false, follow: false },
};

export default function ForgotPasswordPage() {
  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Account"
        title="Forgot password"
        description="If an account exists, a reset path is created. Email delivery is not configured."
      />
      <PageContainer className="max-w-md space-y-6 py-10">
        {isApiConfigured() ? (
          <AuthForm mode="forgot" />
        ) : (
          <p className="text-sm text-muted-foreground">
            NETPULSE_API_BASE_URL is not set. Password reset is unavailable.
          </p>
        )}
      </PageContainer>
    </main>
  );
}
