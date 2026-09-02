import { NextResponse } from "next/server";

import { csrfAllowed } from "@/lib/auth/csrf";
import { ORG_COOKIE, orgCookieOptions, SESSION_COOKIE } from "@/lib/auth/session";
import { getApiBaseUrl } from "@/lib/api/backend";

export const dynamic = "force-dynamic";

export async function POST(request: Request) {
  if (!csrfAllowed(request)) {
    return NextResponse.json(
      { ok: false, error: { code: "forbidden", message: "Request origin was rejected." } },
      { status: 403 }
    );
  }
  const base = getApiBaseUrl();
  if (!base) {
    return NextResponse.json(
      { ok: false, error: { code: "unavailable", message: "NETPULSE_API_BASE_URL is not set." } },
      { status: 503 }
    );
  }
  let orgId = "";
  try {
    const body = (await request.json()) as { orgId?: string };
    orgId = typeof body.orgId === "string" ? body.orgId.trim() : "";
  } catch {
    return NextResponse.json(
      { ok: false, error: { code: "validation_error", message: "Expected JSON." } },
      { status: 400 }
    );
  }
  if (!orgId) {
    return NextResponse.json(
      { ok: false, error: { code: "validation_error", message: "Organization id is required." } },
      { status: 400 }
    );
  }
  const token = cookieValue(request, SESSION_COOKIE);
  if (!token) {
    return NextResponse.json(
      { ok: false, error: { code: "unauthorized", message: "Sign in required." } },
      { status: 401 }
    );
  }
  try {
    const upstream = await fetch(`${base}/v1/orgs/${orgId}`, {
      headers: { Accept: "application/json", Authorization: `Session ${token}` },
      cache: "no-store",
    });
    if (!upstream.ok) {
      return NextResponse.json(
        { ok: false, error: { code: "not_found", message: "organization not found" } },
        { status: 404 }
      );
    }
  } catch {
    return NextResponse.json(
      { ok: false, error: { code: "unavailable", message: "The organization service is unavailable." } },
      { status: 503 }
    );
  }
  const response = NextResponse.json({ ok: true });
  response.cookies.set(ORG_COOKIE, orgId, orgCookieOptions());
  return response;
}

function cookieValue(request: Request, name: string): string | null {
  const header = request.headers.get("cookie");
  if (!header) {
    return null;
  }
  for (const part of header.split(";")) {
    const [key, ...rest] = part.trim().split("=");
    if (key === name) {
      return decodeURIComponent(rest.join("="));
    }
  }
  return null;
}
