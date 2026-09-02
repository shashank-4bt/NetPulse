import type { Metadata } from "next";

import { DevelopmentBanner } from "@/components/public/development-banner";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChangeRoleForm,
  CreateTeamForm,
  DeleteTeamButton,
  InviteMemberForm,
  RemoveMemberButton,
} from "@/features/business/actions";
import { BusinessScreen, BusinessUnavailable, NoOrganization } from "@/features/business/unavailable";
import { hasPermission, roleLabel } from "@/domain/business";
import { getOrgInvites, getOrgMembers, getOrgTeams } from "@/lib/api/business";
import { requireBusiness } from "@/lib/business/page";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Organization team",
  robots: { index: false, follow: false },
};

export default async function BusinessTeamPage() {
  const ctx = await requireBusiness("/business/team");
  if (ctx.unavailable) {
    return <BusinessUnavailable title="Team" description="Organization membership needs the API." />;
  }
  if (!ctx.organization || !ctx.token) {
    return <NoOrganization />;
  }
  const org = ctx.organization;
  const canManage = hasPermission(ctx.permissions, "team.manage");
  const members = await getOrgMembers(ctx.token, org.id);
  const invites = canManage ? await getOrgInvites(ctx.token, org.id) : null;
  const teams = await getOrgTeams(ctx.token, org.id);

  return (
    <BusinessScreen title="Team" description="Invite, remove, and change roles for this organization only.">
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Members</CardTitle>
          </CardHeader>
          <CardContent>
            {!members.ok ? (
              <DevelopmentBanner description={members.message} />
            ) : members.members.length === 0 ? (
              <p className="text-sm text-muted-foreground">No members are stored.</p>
            ) : (
              <ul className="space-y-4 text-sm">
                {members.members.map((item) => (
                  <li key={item.id} className="space-y-2 rounded-md border border-border p-3">
                    <p className="font-medium">{item.displayName || item.email}</p>
                    <p className="text-muted-foreground">
                      {item.email} · {roleLabel(item.role)}
                    </p>
                    {canManage ? (
                      <div className="flex flex-wrap items-end gap-2">
                        <ChangeRoleForm orgId={org.id} memberId={item.id} role={item.role} />
                        <RemoveMemberButton orgId={org.id} memberId={item.id} />
                      </div>
                    ) : null}
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
        {canManage ? (
          <Card>
            <CardHeader>
              <CardTitle>Invite</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <InviteMemberForm orgId={org.id} />
              {!invites || !invites.ok ? (
                <p className="text-sm text-muted-foreground">{invites && !invites.ok ? invites.message : "Invites unavailable."}</p>
              ) : invites.invites.length ? (
                <ul className="space-y-2 text-sm text-muted-foreground">
                  {invites.invites.map((item) => (
                    <li key={item.id}>
                      {item.email} · {roleLabel(item.role)} · {item.summary}
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="text-sm text-muted-foreground">No pending invites are stored.</p>
              )}
            </CardContent>
          </Card>
        ) : null}
        <Card>
          <CardHeader>
            <CardTitle>Teams</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {!teams.ok ? (
              <DevelopmentBanner description={teams.message} />
            ) : teams.teams.length === 0 ? (
              <p className="text-sm text-muted-foreground">No teams are stored.</p>
            ) : (
              <ul className="space-y-3 text-sm">
                {teams.teams.map((item) => (
                  <li key={item.id} className="flex flex-wrap items-start justify-between gap-2">
                    <div>
                      <p className="font-medium">{item.name}</p>
                      <p className="text-muted-foreground">Members listed: {item.memberIds.length}</p>
                    </div>
                    {canManage ? <DeleteTeamButton orgId={org.id} id={item.id} /> : null}
                  </li>
                ))}
              </ul>
            )}
            {canManage ? <CreateTeamForm orgId={org.id} /> : null}
          </CardContent>
        </Card>
      </div>
    </BusinessScreen>
  );
}
