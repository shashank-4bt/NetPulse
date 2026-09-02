import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CreateOrgKeyForm, RenameOrgForm, RevokeOrgKeyButton, RotateOrgKeyButton } from "@/features/business/actions";
import { BusinessScreen, BusinessUnavailable, NoOrganization } from "@/features/business/unavailable";
import { hasPermission, roleLabel } from "@/domain/business";
import { getOrgAudit, getOrgBilling, getOrgKeys } from "@/lib/api/business";
import { requireBusiness } from "@/lib/business/page";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Organization settings",
  robots: { index: false, follow: false },
};

export default async function BusinessSettingsPage() {
  const ctx = await requireBusiness("/business/settings");
  if (ctx.unavailable) {
    return <BusinessUnavailable title="Settings" description="Organization settings need the API." />;
  }
  if (!ctx.organization || !ctx.token) {
    return <NoOrganization />;
  }
  const org = ctx.organization;
  const canTeam = hasPermission(ctx.permissions, "team.manage");
  const canBilling = hasPermission(ctx.permissions, "billing.manage");
  const canKeys = hasPermission(ctx.permissions, "api.manage");
  const canAudit = hasPermission(ctx.permissions, "audit.read");
  const billing = canBilling ? await getOrgBilling(ctx.token, org.id) : null;
  const keys = canKeys ? await getOrgKeys(ctx.token, org.id) : null;
  const audit = canAudit ? await getOrgAudit(ctx.token, org.id) : null;

  return (
    <BusinessScreen
      title="Settings"
      description="Organization name, billing, API keys, and audit events for the selected membership."
    >
      <p className="text-sm text-muted-foreground">
        Your role: {roleLabel(org.role)}. Permissions are enforced by the API, not this page.
      </p>
      {canTeam ? (
        <Card>
          <CardHeader>
            <CardTitle>Organization name</CardTitle>
          </CardHeader>
          <CardContent>
            <RenameOrgForm orgId={org.id} name={org.name} />
          </CardContent>
        </Card>
      ) : null}
      <Card>
        <CardHeader>
          <CardTitle>Billing</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2 text-sm text-muted-foreground">
          {!canBilling ? (
            <p>Missing permission. Billing is not shown.</p>
          ) : !billing || !billing.ok ? (
            <DevelopmentBanner description={billing && !billing.ok ? billing.message : "Billing unavailable."} />
          ) : (
            <>
              <p>{billing.billing.summary}</p>
              <p>Has billing account: {billing.billing.hasAccount ? "yes" : "no"}.</p>
              <p>Invoices stored: {billing.billing.invoices.length}.</p>
            </>
          )}
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Organization API keys</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {!canKeys ? (
            <p className="text-sm text-muted-foreground">Missing permission. Keys are not shown.</p>
          ) : !keys || !keys.ok ? (
            <DevelopmentBanner description={keys && !keys.ok ? keys.message : "Keys unavailable."} />
          ) : (
            <>
              {keys.keys.length === 0 ? (
                <p className="text-sm text-muted-foreground">No organization keys are stored.</p>
              ) : (
                <ul className="space-y-3 text-sm">
                  {keys.keys.map((item) => (
                    <li key={item.id} className="flex flex-wrap items-start justify-between gap-2">
                      <div>
                        <p className="font-medium">
                          {item.name} · {item.prefix}…{item.last4}
                        </p>
                        <p className="text-muted-foreground">
                          {item.revoked ? "revoked" : "active"} · {item.scopes.join(", ") || "no scopes"}
                        </p>
                      </div>
                      {!item.revoked ? (
                        <div className="flex gap-2">
                          <RotateOrgKeyButton orgId={org.id} id={item.id} />
                          <RevokeOrgKeyButton orgId={org.id} id={item.id} />
                        </div>
                      ) : null}
                    </li>
                  ))}
                </ul>
              )}
              <CreateOrgKeyForm orgId={org.id} />
            </>
          )}
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Audit</CardTitle>
        </CardHeader>
        <CardContent>
          {!canAudit ? (
            <p className="text-sm text-muted-foreground">Missing permission. Audit events are not shown.</p>
          ) : !audit || !audit.ok ? (
            <DevelopmentBanner description={audit && !audit.ok ? audit.message : "Audit unavailable."} />
          ) : audit.events.length === 0 ? (
            <p className="text-sm text-muted-foreground">No audit events are stored.</p>
          ) : (
            <ul className="space-y-2 text-sm">
              {audit.events.map((item) => (
                <li key={item.id}>
                  <p className="font-medium">{item.kind}</p>
                  <p className="text-muted-foreground">
                    {item.at} · {item.summary}
                  </p>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </BusinessScreen>
  );
}
