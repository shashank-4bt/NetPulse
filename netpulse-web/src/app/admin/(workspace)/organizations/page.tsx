import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { AdminForbidden, AdminScreen, AdminUnavailable } from "@/features/admin/screen";
import { requireAdmin } from "@/lib/admin/page";
import { getAdminOrganizations } from "@/lib/api/admin";

export const dynamic = "force-dynamic";

export const metadata: Metadata = { title: "Organizations", robots: { index: false, follow: false } };

export default async function AdminOrganizationsPage() {
  const ctx = await requireAdmin("/admin/organizations");
  if (ctx.unavailable) {
    return <AdminUnavailable title="Organizations" description="Organization records stay unavailable until the API is connected." />;
  }
  if (ctx.forbidden) {
    return <AdminForbidden />;
  }
  const loaded = await getAdminOrganizations(ctx.token);
  return (
    <AdminScreen title="Organizations" description="Stored organizations only. Membership secrets are not listed.">
      {!loaded.ok ? (
        <DevelopmentBanner description={loaded.message} />
      ) : loaded.organizations.length === 0 ? (
        <p className="text-sm text-muted-foreground">No organizations are stored.</p>
      ) : (
        <ul className="space-y-3 text-sm">
          {loaded.organizations.map((item) => (
            <li key={item.id}>
              <p className="font-medium">{item.name}</p>
              <p className="text-muted-foreground">{item.id} · {item.createdAt}</p>
            </li>
          ))}
        </ul>
      )}
    </AdminScreen>
  );
}
