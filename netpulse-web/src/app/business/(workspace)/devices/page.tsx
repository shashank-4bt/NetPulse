import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CreateDeviceForm, DeleteDeviceButton } from "@/features/business/actions";
import { BusinessScreen, BusinessUnavailable, NoOrganization } from "@/features/business/unavailable";
import { hasPermission } from "@/domain/business";
import { getOrgDevices } from "@/lib/api/business";
import { requireBusiness } from "@/lib/business/page";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Organization devices",
  robots: { index: false, follow: false },
};

export default async function BusinessDevicesPage() {
  const ctx = await requireBusiness("/business/devices");
  if (ctx.unavailable) {
    return <BusinessUnavailable title="Devices" description="Organization devices need the API." />;
  }
  if (!ctx.organization || !ctx.token) {
    return <NoOrganization />;
  }
  const org = ctx.organization;
  const loaded = await getOrgDevices(ctx.token, org.id);
  const canWrite = hasPermission(ctx.permissions, "monitor.create") || hasPermission(ctx.permissions, "team.manage");
  const canDelete = hasPermission(ctx.permissions, "monitor.delete") || hasPermission(ctx.permissions, "team.manage");

  return (
    <BusinessScreen title="Devices" description="Stored device labels for this organization. These are not discovered hosts.">
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Stored devices</CardTitle>
          </CardHeader>
          <CardContent>
            {!loaded.ok ? (
              <DevelopmentBanner description={loaded.message} />
            ) : loaded.devices.length === 0 ? (
              <p className="text-sm text-muted-foreground">No devices are stored.</p>
            ) : (
              <ul className="space-y-3 text-sm">
                {loaded.devices.map((item) => (
                  <li key={item.id} className="flex flex-wrap items-start justify-between gap-2">
                    <div>
                      <p className="font-medium">{item.name}</p>
                      <p className="text-muted-foreground">
                        {item.label || "no label"} · {item.region || "no region"}
                      </p>
                      <p className="text-muted-foreground">{item.summary}</p>
                    </div>
                    {canDelete ? <DeleteDeviceButton orgId={org.id} id={item.id} /> : null}
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
        {canWrite ? (
          <Card>
            <CardHeader>
              <CardTitle>Add device</CardTitle>
            </CardHeader>
            <CardContent>
              <CreateDeviceForm orgId={org.id} />
            </CardContent>
          </Card>
        ) : null}
      </div>
    </BusinessScreen>
  );
}
