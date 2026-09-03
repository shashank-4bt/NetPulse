import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { AdminForbidden, AdminScreen, AdminUnavailable } from "@/features/admin/screen";
import { requireAdmin } from "@/lib/admin/page";
import { getAdminUsers } from "@/lib/api/admin";

export const dynamic = "force-dynamic";

export const metadata: Metadata = { title: "Users", robots: { index: false, follow: false } };

export default async function AdminUsersPage() {
  const ctx = await requireAdmin("/admin/users");
  if (ctx.unavailable) {
    return <AdminUnavailable title="Users" description="Account records stay unavailable until the API is connected." />;
  }
  if (ctx.forbidden) {
    return <AdminForbidden />;
  }
  const loaded = await getAdminUsers(ctx.token);
  return (
    <AdminScreen title="Users" description="Password hashes and session secrets are never listed.">
      {!loaded.ok ? (
        <DevelopmentBanner description={loaded.message} />
      ) : loaded.users.length === 0 ? (
        <p className="text-sm text-muted-foreground">No accounts are stored.</p>
      ) : (
        <ul className="space-y-3 text-sm">
          {loaded.users.map((item) => (
            <li key={item.id}>
              <p className="font-medium">{item.displayName || item.email}</p>
              <p className="text-muted-foreground">
                {item.email} · {item.emailVerified ? "verified" : "unverified"} · {item.createdAt}
              </p>
              <p className="text-muted-foreground">{item.summary}</p>
            </li>
          ))}
        </ul>
      )}
    </AdminScreen>
  );
}
