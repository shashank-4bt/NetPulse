import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { AdminForbidden, AdminScreen, AdminUnavailable } from "@/features/admin/screen";
import { requireAdmin } from "@/lib/admin/page";
import { getAdminDiagnostics } from "@/lib/api/admin";

export const dynamic = "force-dynamic";

export const metadata: Metadata = { title: "Diagnostics", robots: { index: false, follow: false } };

export default async function AdminDiagnosticsPage() {
  const ctx = await requireAdmin("/admin/diagnostics");
  if (ctx.unavailable) {
    return <AdminUnavailable title="Diagnostics" description="Stored diagnoses stay unavailable until the API is connected." />;
  }
  if (ctx.forbidden) {
    return <AdminForbidden />;
  }
  const loaded = await getAdminDiagnostics(ctx.token);
  return (
    <AdminScreen title="Diagnostics" description="Stored diagnosis outcomes from worker vantages. Share tokens are not listed.">
      {!loaded.ok ? (
        <DevelopmentBanner description={loaded.message} />
      ) : loaded.diagnoses.length === 0 ? (
        <p className="text-sm text-muted-foreground">No diagnoses are stored.</p>
      ) : (
        <ul className="space-y-3 text-sm">
          {loaded.diagnoses.map((item) => (
            <li key={item.id}>
              <p className="font-medium">
                {item.target} · {item.status}
              </p>
              <p className="text-muted-foreground">{item.id} · {item.createdAt}</p>
              <p className="text-muted-foreground">{item.summary}</p>
            </li>
          ))}
        </ul>
      )}
    </AdminScreen>
  );
}
