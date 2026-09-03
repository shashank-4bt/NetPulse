import { proxyAccount } from "@/lib/auth/bff";

export const dynamic = "force-dynamic";

type RouteCtx = { params: Promise<{ path: string[] }> };

async function handle(request: Request, ctx: RouteCtx) {
  return proxyAccount(request, "admin", await ctx.params);
}

export const GET = handle;
export const POST = handle;
export const PUT = handle;
export const PATCH = handle;
export const DELETE = handle;
