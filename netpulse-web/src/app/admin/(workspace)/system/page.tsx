import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CreateFlagForm, UpdateConfigForm } from "@/features/admin/actions";
import { SystemHealth } from "@/features/admin/health";
import { AdminForbidden, AdminScreen, AdminUnavailable } from "@/features/admin/screen";
import { requireAdmin } from "@/lib/admin/page";
import { getAdminConfig, getAdminFlags, getAdminSystem } from "@/lib/api/admin";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "System",
  robots: { index: false, follow: false },
};

export default async function AdminSystemPage() {
  const ctx = await requireAdmin("/admin/system");
  if (ctx.unavailable) {
    return <AdminUnavailable title="System" description="Operational health stays unavailable until the API is connected." />;
  }
  if (ctx.forbidden) {
    return <AdminForbidden />;
  }
  const [system, flags, config] = await Promise.all([
    getAdminSystem(ctx.token),
    getAdminFlags(ctx.token),
    getAdminConfig(ctx.token),
  ]);
  return (
    <AdminScreen
      title="System"
      description="Remote configuration and feature flags are stored settings. They are not live traffic percentages."
    >
      {!system.ok ? <DevelopmentBanner description={system.message} /> : <SystemHealth system={system.system} />}
      <Card>
        <CardHeader>
          <CardTitle>Feature flags</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {!flags.ok ? (
            <DevelopmentBanner description={flags.message} />
          ) : flags.flags.length === 0 ? (
            <p className="text-sm text-muted-foreground">No feature flags are stored.</p>
          ) : (
            <ul className="space-y-3 text-sm">
              {flags.flags.map((item) => (
                <li key={item.id}>
                  <p className="font-medium">
                    {item.name} · {item.enabled ? "enabled" : "disabled"} · {item.percentage}%
                  </p>
                  <p className="text-muted-foreground">
                    Environment {item.environment || "any"} · users {item.userIds.length} · orgs {item.orgIds.length}
                  </p>
                  <p className="text-muted-foreground">{item.summary}</p>
                </li>
              ))}
            </ul>
          )}
          <CreateFlagForm />
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Remote configuration</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {!config.ok ? (
            <DevelopmentBanner description={config.message} />
          ) : (
            <ul className="space-y-3 text-sm">
              {config.entries.map((item) => (
                <li key={item.key}>
                  <p className="font-medium">
                    {item.key} = {item.value}
                  </p>
                  <p className="text-muted-foreground">{item.summary}</p>
                </li>
              ))}
            </ul>
          )}
          <UpdateConfigForm />
        </CardContent>
      </Card>
    </AdminScreen>
  );
}
