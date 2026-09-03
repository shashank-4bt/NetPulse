import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { AdminForbidden, AdminScreen, AdminUnavailable } from "@/features/admin/screen";
import { requireAdmin } from "@/lib/admin/page";
import { getAdminServices } from "@/lib/api/admin";

export const dynamic = "force-dynamic";

export const metadata: Metadata = { title: "Services", robots: { index: false, follow: false } };

export default async function AdminServicesPage() {
  const ctx = await requireAdmin("/admin/services");
  if (ctx.unavailable) {
    return <AdminUnavailable title="Services" description="The editorial catalog stays unavailable until the API is connected." />;
  }
  if (ctx.forbidden) {
    return <AdminForbidden />;
  }
  const loaded = await getAdminServices(ctx.token);
  return (
    <AdminScreen title="Services" description="Editorial catalog entries. This is not a live health feed.">
      {!loaded.ok ? (
        <DevelopmentBanner description={loaded.message} />
      ) : loaded.services.length === 0 ? (
        <p className="text-sm text-muted-foreground">No catalog entries are stored.</p>
      ) : (
        <ul className="space-y-3 text-sm">
          {loaded.services.map((item) => (
            <li key={item.slug}>
              <p className="font-medium">{item.name}</p>
              <p className="text-muted-foreground">{item.category} · {item.summary}</p>
            </li>
          ))}
        </ul>
      )}
    </AdminScreen>
  );
}
