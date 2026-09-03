import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { AdminForbidden, AdminScreen, AdminUnavailable } from "@/features/admin/screen";
import { requireAdmin } from "@/lib/admin/page";
import { getAdminAudit } from "@/lib/api/admin";

export const dynamic = "force-dynamic";

export const metadata: Metadata = { title: "Audit", robots: { index: false, follow: false } };

export default async function AdminAuditPage() {
  const ctx = await requireAdmin("/admin/audit");
  if (ctx.unavailable) {
    return <AdminUnavailable title="Audit" description="Operator audit events stay unavailable until the API is connected." />;
  }
  if (ctx.forbidden) {
    return <AdminForbidden />;
  }
  const loaded = await getAdminAudit(ctx.token);
  return (
    <AdminScreen title="Audit" description="Actor, action, resource, timestamp, and result for operator actions.">
      {!loaded.ok ? (
        <DevelopmentBanner description={loaded.message} />
      ) : loaded.events.length === 0 ? (
        <p className="text-sm text-muted-foreground">No operator audit events are stored.</p>
      ) : (
        <ul className="space-y-3 text-sm">
          {loaded.events.map((item) => (
            <li key={item.id}>
              <p className="font-medium">
                {item.action} · {item.result}
              </p>
              <p className="text-muted-foreground">
                {item.at} · {item.actorId} · {item.resource}
              </p>
              <p className="text-muted-foreground">{item.summary}</p>
            </li>
          ))}
        </ul>
      )}
    </AdminScreen>
  );
}
