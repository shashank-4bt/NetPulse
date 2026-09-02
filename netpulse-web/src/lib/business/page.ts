import { requireAccount } from "@/lib/account/guard";
import { loadOrgContext, type OrgContext } from "@/lib/business/org";
import { readSessionToken } from "@/lib/auth/session";

export async function requireBusiness(path: string): Promise<
  | { unavailable: true }
  | ({ unavailable: false; token: string | null } & OrgContext)
> {
  const gate = await requireAccount(path);
  if (gate.unavailable) {
    return { unavailable: true };
  }
  const token = await readSessionToken();
  const ctx = token ? await loadOrgContext(token) : { organizations: [], organization: null, permissions: [] };
  return { unavailable: false, token, ...ctx };
}
