import { proxyAccount } from "@/lib/auth/bff";

export const dynamic = "force-dynamic";

function handle(request: Request) {
  return proxyAccount(request, "orgs", { path: [] });
}

export const GET = handle;
export const POST = handle;
