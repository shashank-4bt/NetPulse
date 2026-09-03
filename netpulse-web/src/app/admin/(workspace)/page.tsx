import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { AdminForbidden, AdminScreen, AdminUnavailable } from "@/features/admin/screen";
import { SystemHealth } from "@/features/admin/health";
import { requireAdmin } from "@/lib/admin/page";
import { getAdminSystem } from "@/lib/api/admin";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Operations",
  robots: { index: false, follow: false },
};

export default async function AdminOverviewPage() {
  const ctx = await requireAdmin("/admin");
  if (ctx.unavailable) {
    return (
      <AdminUnavailable
        title="Operations"
        description="API health, worker health, and error rates stay unavailable until the API is connected."
      />
    );
  }
  if (ctx.forbidden) {
    return <AdminForbidden />;
  }
  const loaded = await getAdminSystem(ctx.token);
  return (
    <AdminScreen
      title="Operations overview"
      description="Figures use this API process, stored worker heartbeats, queue depth, and stored measurements. Empty series stay unmeasured."
    >
      {!loaded.ok ? (
        <DevelopmentBanner description={loaded.message} />
      ) : (
        <SystemHealth system={loaded.system} />
      )}
    </AdminScreen>
  );
}
