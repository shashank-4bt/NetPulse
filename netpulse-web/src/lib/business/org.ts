import type { OrganizationView } from "@/domain/business";
import { getOrg, listOrgs } from "@/lib/api/business";
import { readOrgId } from "@/lib/auth/session";

export type OrgContext = {
  organizations: OrganizationView[];
  organization: OrganizationView | null;
  permissions: string[];
};

export async function loadOrgContext(session: string): Promise<OrgContext> {
  const listed = await listOrgs(session);
  if (!listed.ok) {
    return { organizations: [], organization: null, permissions: [] };
  }
  const selected = await readOrgId();
  const match = listed.organizations.find((item) => item.id === selected) ?? listed.organizations[0] ?? null;
  if (!match) {
    return { organizations: listed.organizations, organization: null, permissions: [] };
  }
  const detail = await getOrg(session, match.id);
  return {
    organizations: listed.organizations,
    organization: detail.ok ? detail.organization : match,
    permissions: detail.ok ? detail.permissions : [],
  };
}
