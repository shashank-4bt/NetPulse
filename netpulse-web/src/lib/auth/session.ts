import { cookies } from "next/headers";

import type { AccountUser } from "@/domain/account";
import { getAuthMe, isApiConfigured } from "@/lib/api/account";

export const SESSION_COOKIE = "np_session";
export const ORG_COOKIE = "np_org";
export const SESSION_MAX_AGE = 60 * 60 * 24 * 7;

export function sessionCookieOptions() {
  return {
    httpOnly: true,
    sameSite: "lax" as const,
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: SESSION_MAX_AGE,
  };
}

export function orgCookieOptions() {
  return {
    httpOnly: true,
    sameSite: "lax" as const,
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: SESSION_MAX_AGE,
  };
}

export async function readOrgId(): Promise<string | null> {
  try {
    const jar = await cookies();
    const value = jar.get(ORG_COOKIE)?.value?.trim();
    return value || null;
  } catch {
    return null;
  }
}

export async function readSessionToken(): Promise<string | null> {
  try {
    const jar = await cookies();
    const value = jar.get(SESSION_COOKIE)?.value?.trim();
    return value || null;
  } catch {
    return null;
  }
}

export async function getCurrentUser(): Promise<AccountUser | null> {
  const token = await readSessionToken();
  if (!token || !isApiConfigured()) {
    return null;
  }
  const result = await getAuthMe(token);
  return result.ok ? result.user : null;
}
