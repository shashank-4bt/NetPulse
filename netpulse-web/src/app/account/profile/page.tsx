import type { Metadata } from "next";

import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { LogoutButton, ProfileForm } from "@/features/account/account-actions";
import { requireAccount } from "@/lib/account/guard";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Profile",
  robots: { index: false, follow: false },
};

export default async function ProfilePage() {
  const { user, unavailable } = await requireAccount("/account/profile");
  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Account"
        title="Profile"
        description="Email is the account identifier. Display name is the only profile field you can edit here."
      />
      <PageContainer className="max-w-xl space-y-6 py-10">
        {unavailable || !user ? (
          <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set. Profile editing is unavailable." />
        ) : (
          <>
            <p className="text-sm text-muted-foreground">
              {user.email}
              {user.emailVerified ? " · email verified" : " · email not verified"}
            </p>
            <ProfileForm displayName={user.displayName} />
            <LogoutButton />
          </>
        )}
      </PageContainer>
    </main>
  );
}
