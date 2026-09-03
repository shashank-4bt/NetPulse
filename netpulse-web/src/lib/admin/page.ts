import { requireAccount } from "@/lib/account/guard";
import { getAdminMe } from "@/lib/api/admin";
import { readSessionToken } from "@/lib/auth/session";
import type { OperatorView } from "@/domain/admin";

export async function requireAdmin(path: string): Promise<
  | { unavailable: true; forbidden: false; token: null; operator: null }
  | { unavailable: false; forbidden: true; token: string | null; operator: null }
  | { unavailable: false; forbidden: false; token: string; operator: OperatorView }
> {
  const gate = await requireAccount(path);
  if (gate.unavailable) {
    return { unavailable: true, forbidden: false, token: null, operator: null };
  }
  const token = await readSessionToken();
  if (!token) {
    return { unavailable: true, forbidden: false, token: null, operator: null };
  }
  const me = await getAdminMe(token);
  if (!me.ok) {
    if (me.status === 404) {
      return { unavailable: false, forbidden: true, token, operator: null };
    }
    return { unavailable: true, forbidden: false, token: null, operator: null };
  }
  return { unavailable: false, forbidden: false, token, operator: me.operator };
}
