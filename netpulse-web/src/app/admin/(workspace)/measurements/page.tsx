import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { AdminForbidden, AdminScreen, AdminUnavailable } from "@/features/admin/screen";
import { requireAdmin } from "@/lib/admin/page";
import { getAdminMeasurements } from "@/lib/api/admin";

export const dynamic = "force-dynamic";

export const metadata: Metadata = { title: "Measurements", robots: { index: false, follow: false } };

export default async function AdminMeasurementsPage() {
  const ctx = await requireAdmin("/admin/measurements");
  if (ctx.unavailable) {
    return <AdminUnavailable title="Measurements" description="Stored measurements stay unavailable until the API is connected." />;
  }
  if (ctx.forbidden) {
    return <AdminForbidden />;
  }
  const loaded = await getAdminMeasurements(ctx.token);
  return (
    <AdminScreen title="Measurements" description="Worker vantage observations only. Unmeasured rows stay unlabeled.">
      {!loaded.ok ? (
        <DevelopmentBanner description={loaded.message} />
      ) : loaded.measurements.length === 0 ? (
        <p className="text-sm text-muted-foreground">No measurements are stored.</p>
      ) : (
        <ul className="space-y-3 text-sm">
          {loaded.measurements.map((item, index) => (
            <li key={`${item.diagnosisId}-${item.key}-${index}`}>
              <p className="font-medium">
                {item.label} · {item.measured ? "measured" : "unmeasured"}
              </p>
              <p className="text-muted-foreground">{item.diagnosisId} · {item.summary}</p>
            </li>
          ))}
        </ul>
      )}
    </AdminScreen>
  );
}
