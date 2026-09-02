import { redirect } from "next/navigation";

import type { AccountUser } from "@/domain/account";
import { isApiConfigured } from "@/lib/api/backend";
import { getCurrentUser } from "@/lib/auth/session";

export async function requireAccount(nextPath: string): Promise<{
  user: AccountUser | null;
  unavailable: boolean;
}> {
  if (!isApiConfigured()) {
    return { user: null, unavailable: true };
  }
  const user = await getCurrentUser();
  if (!user) {
    redirect(`/login?next=${encodeURIComponent(nextPath)}`);
  }
  return { user, unavailable: false };
}

export async function redirectIfSignedIn() {
  if (!isApiConfigured()) {
    return;
  }
  const user = await getCurrentUser();
  if (user) {
    redirect("/dashboard");
  }
}
