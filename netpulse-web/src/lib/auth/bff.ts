import { NextResponse } from "next/server";

import { csrfAllowed } from "@/lib/auth/csrf";
import { SESSION_COOKIE, sessionCookieOptions } from "@/lib/auth/session";
import { getApiBaseUrl } from "@/lib/api/backend";

type PathParams = { path?: string[] };

export async function proxyAccount(
  request: Request,
  prefix: "auth" | "me" | "dev",
  params: PathParams
): Promise<NextResponse> {
  if (!csrfAllowed(request)) {
    return NextResponse.json(
      { ok: false, error: { code: "forbidden", message: "Request origin was rejected." } },
      { status: 403 }
    );
  }

  const base = getApiBaseUrl();
  if (!base) {
    return NextResponse.json(
      {
        ok: false,
        error: {
          code: "unavailable",
          message: "NETPULSE_API_BASE_URL is not set. Accounts are unavailable.",
        },
      },
      { status: 503 }
    );
  }

  const suffix = (params.path ?? []).join("/");
  const incoming = new URL(request.url);
  const target = `${base}/v1/${prefix}/${suffix}${incoming.search}`;
  const token = cookieValue(request, SESSION_COOKIE);
  const headers = new Headers();
  headers.set("Accept", "application/json");
  if (request.headers.get("content-type")) {
    headers.set("Content-Type", request.headers.get("content-type") ?? "application/json");
  }
  if (token) {
    headers.set("Authorization", `Session ${token}`);
  }
  const forwarded = request.headers.get("x-forwarded-for");
  if (forwarded) {
    headers.set("X-Forwarded-For", forwarded);
  }
  const ua = request.headers.get("user-agent");
  if (ua) {
    headers.set("User-Agent", ua);
  }

  let upstream: Response;
  try {
    upstream = await fetch(target, {
      method: request.method,
      headers,
      body: canHaveBody(request.method) ? await request.text() : undefined,
      cache: "no-store",
    });
  } catch {
    return NextResponse.json(
      {
        ok: false,
        error: {
          code: "unavailable",
          message: "The account service is unavailable. No fallback account was created.",
        },
      },
      { status: 503 }
    );
  }

  let payload: Record<string, unknown> = {};
  try {
    payload = (await upstream.json()) as Record<string, unknown>;
  } catch {
    payload = {};
  }

  const sessionToken =
    typeof payload.sessionToken === "string" ? payload.sessionToken : "";
  delete payload.sessionToken;
  if (process.env.NODE_ENV === "production" && payload.auth && typeof payload.auth === "object") {
    delete (payload.auth as { devToken?: string }).devToken;
  }

  const response = NextResponse.json(payload, { status: upstream.status });
  const path = `/${suffix}`;
  if (sessionToken && (path === "/register" || path === "/login")) {
    response.cookies.set(SESSION_COOKIE, sessionToken, sessionCookieOptions());
  }
  if (path === "/logout" || path === "/deletion") {
    response.cookies.set(SESSION_COOKIE, "", { ...sessionCookieOptions(), maxAge: 0 });
  }
  return response;
}

function canHaveBody(method: string): boolean {
  return method !== "GET" && method !== "HEAD";
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
