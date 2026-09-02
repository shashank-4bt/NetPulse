import type { Metadata } from "next";

import { EmptyState } from "@/components/feedback/empty-state";
import { PageContainer } from "@/components/layout/page-container";
import { DevelopmentBanner } from "@/components/public/development-banner";
import { PageHero } from "@/components/public/page-hero";
import { RevokeSessionButton } from "@/features/account/account-actions";
import { requireAccount } from "@/lib/account/guard";
import { getDevices } from "@/lib/api/account";
import { readSessionToken } from "@/lib/auth/session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Devices",
  robots: { index: false, follow: false },
};

export default async function DevicesPage() {
  const { unavailable } = await requireAccount("/account/devices");
  if (unavailable) {
    return (
      <main id="main-content" className="flex-1">
        <PageHero eyebrow="Account" title="Devices" description="Device list needs the API." />
        <PageContainer className="py-10">
          <DevelopmentBanner description="NETPULSE_API_BASE_URL is not set." />
        </PageContainer>
      </main>
    );
  }
  const token = await readSessionToken();
  const loaded = token ? await getDevices(token) : null;

  return (
    <main id="main-content" className="flex-1">
      <PageHero
        eyebrow="Account"
        title="Devices"
        description="These are sign-in sessions, not inventoried hardware. IPs are coarsened. NetPulse does not collect a device inventory from the browser."
      />
      <PageContainer className="max-w-2xl space-y-4 py-10">
        {!loaded || !loaded.ok ? (
          <EmptyState title="Devices unavailable" description={loaded && !loaded.ok ? loaded.message : ""} />
        ) : loaded.devices.length === 0 ? (
          <EmptyState title="No sessions" description="No active sessions are stored." />
        ) : (
          <ul className="space-y-3">
            {loaded.devices.map((device) => (
              <li key={device.id} className="rounded-lg border border-border p-3 text-sm">
                <p className="font-medium">
                  {device.label}
                  {device.current ? " · this session" : ""}
                </p>
                <p className="text-muted-foreground">
                  {device.kind} · {device.ip || "IP not stored"} · {device.lastSeenAt}
                </p>
                {device.current ? null : <RevokeSessionButton sessionId={device.id} />}
              </li>
            ))}
          </ul>
        )}
      </PageContainer>
    </main>
  );
}
