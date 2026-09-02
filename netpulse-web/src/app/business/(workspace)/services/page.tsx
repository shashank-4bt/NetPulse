import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CreateOrgServiceForm, DeleteOrgServiceButton } from "@/features/business/actions";
import { BusinessScreen, BusinessUnavailable, NoOrganization } from "@/features/business/unavailable";
import { hasPermission } from "@/domain/business";
import { getOrgServices } from "@/lib/api/business";
import { requireBusiness } from "@/lib/business/page";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Organization services",
  robots: { index: false, follow: false },
};

export default async function BusinessServicesPage() {
  const ctx = await requireBusiness("/business/services");
  if (ctx.unavailable) {
    return <BusinessUnavailable title="Services" description="Organization services need the API." />;
  }
  if (!ctx.organization || !ctx.token) {
    return <NoOrganization />;
  }
  const org = ctx.organization;
  const loaded = await getOrgServices(ctx.token, org.id);
  const canWrite = hasPermission(ctx.permissions, "monitor.create") || hasPermission(ctx.permissions, "team.manage");
  const canDelete = hasPermission(ctx.permissions, "monitor.delete") || hasPermission(ctx.permissions, "team.manage");

  return (
    <BusinessScreen title="Services" description="Stored service endpoints for this organization. This is not the public catalog.">
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Stored services</CardTitle>
          </CardHeader>
          <CardContent>
            {!loaded.ok ? (
              <DevelopmentBanner description={loaded.message} />
            ) : loaded.services.length === 0 ? (
              <p className="text-sm text-muted-foreground">No services are stored.</p>
            ) : (
              <ul className="space-y-3 text-sm">
                {loaded.services.map((item) => (
                  <li key={item.id} className="flex flex-wrap items-start justify-between gap-2">
                    <div>
                      <p className="font-medium">{item.name}</p>
                      <p className="text-muted-foreground">
                        {item.slug || "no slug"} · {item.endpoint || "no endpoint"}
                      </p>
                      <p className="text-muted-foreground">{item.summary}</p>
                    </div>
                    {canDelete ? <DeleteOrgServiceButton orgId={org.id} id={item.id} /> : null}
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
        {canWrite ? (
          <Card>
            <CardHeader>
              <CardTitle>Add service</CardTitle>
            </CardHeader>
            <CardContent>
              <CreateOrgServiceForm orgId={org.id} />
            </CardContent>
          </Card>
        ) : null}
      </div>
    </BusinessScreen>
  );
}
