import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CreateNetworkForm, DeleteNetworkButton } from "@/features/business/actions";
import { BusinessScreen, BusinessUnavailable, NoOrganization } from "@/features/business/unavailable";
import { hasPermission } from "@/domain/business";
import { getOrgNetworks } from "@/lib/api/business";
import { requireBusiness } from "@/lib/business/page";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Organization networks",
  robots: { index: false, follow: false },
};

export default async function BusinessNetworksPage() {
  const ctx = await requireBusiness("/business/networks");
  if (ctx.unavailable) {
    return <BusinessUnavailable title="Networks" description="Organization networks need the API." />;
  }
  if (!ctx.organization || !ctx.token) {
    return <NoOrganization />;
  }
  const org = ctx.organization;
  const loaded = await getOrgNetworks(ctx.token, org.id);
  const canWrite = hasPermission(ctx.permissions, "monitor.create") || hasPermission(ctx.permissions, "team.manage");
  const canDelete = hasPermission(ctx.permissions, "monitor.delete") || hasPermission(ctx.permissions, "team.manage");

  return (
    <BusinessScreen title="Networks" description="Stored network and ASN records for this organization.">
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Stored networks</CardTitle>
          </CardHeader>
          <CardContent>
            {!loaded.ok ? (
              <DevelopmentBanner description={loaded.message} />
            ) : loaded.networks.length === 0 ? (
              <p className="text-sm text-muted-foreground">No networks are stored.</p>
            ) : (
              <ul className="space-y-3 text-sm">
                {loaded.networks.map((item) => (
                  <li key={item.id} className="flex flex-wrap items-start justify-between gap-2">
                    <div>
                      <p className="font-medium">{item.name}</p>
                      <p className="text-muted-foreground">
                        {item.asn || "no ASN"} · {item.region || "no region"}
                      </p>
                      <p className="text-muted-foreground">{item.summary}</p>
                    </div>
                    {canDelete ? <DeleteNetworkButton orgId={org.id} id={item.id} /> : null}
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
        {canWrite ? (
          <Card>
            <CardHeader>
              <CardTitle>Add network</CardTitle>
            </CardHeader>
            <CardContent>
              <CreateNetworkForm orgId={org.id} />
            </CardContent>
          </Card>
        ) : null}
      </div>
    </BusinessScreen>
  );
}
