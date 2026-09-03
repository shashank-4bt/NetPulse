import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { AdminForbidden, AdminScreen, AdminUnavailable } from "@/features/admin/screen";
import { requireAdmin } from "@/lib/admin/page";
import { getAdminAbuse } from "@/lib/api/admin";

export const dynamic = "force-dynamic";

export const metadata: Metadata = { title: "Abuse", robots: { index: false, follow: false } };

export default async function AdminAbusePage() {
  const ctx = await requireAdmin("/admin/abuse");
  if (ctx.unavailable) {
    return <AdminUnavailable title="Abuse" description="Rate-limit and SSRF events stay unavailable until the API is connected." />;
  }
  if (ctx.forbidden) {
    return <AdminForbidden />;
  }
  const loaded = await getAdminAbuse(ctx.token);
  return (
    <AdminScreen title="Abuse" description="Stored rate-limit events, SSRF attempts, and API-key abuse. Volume is not estimated.">
      {!loaded.ok ? (
        <DevelopmentBanner description={loaded.message} />
      ) : loaded.events.length === 0 ? (
        <p className="text-sm text-muted-foreground">No abuse events are stored.</p>
      ) : (
        <ul className="space-y-3 text-sm">
          {loaded.events.map((item) => (
            <li key={item.id}>
              <p className="font-medium">
                {item.kind} · {item.result}
              </p>
              <p className="text-muted-foreground">
                {item.at} · {item.resource} · {item.ip || item.actor}
              </p>
              <p className="text-muted-foreground">{item.summary}</p>
            </li>
          ))}
        </ul>
      )}
    </AdminScreen>
  );
}
