import { NextResponse } from "next/server";

import { emptyMapAggregates } from "@/domain/map";
import { parseMapQuery } from "@/features/map/query";
import { roundViewport } from "@/features/map/privacy";
import { getBackendMapAggregates, isApiConfigured } from "@/lib/api/backend";

export const dynamic = "force-dynamic";

export async function GET(request: Request) {
  const url = new URL(request.url);
  const params = Object.fromEntries(url.searchParams.entries());
  const layers = url.searchParams.getAll("layers");
  const query = parseMapQuery({
    ...params,
    layers: layers.length ? layers : params.layers,
  });
  const viewport = parseViewport(url);

  if (!isApiConfigured()) {
    const reason =
      "NETPULSE_API_BASE_URL is not set. Viewport queries will not invent map cells.";
    return NextResponse.json(
      {
        ok: false,
        error: { code: "unavailable", message: reason },
        map: emptyMapAggregates(reason),
      },
      { status: 503 }
    );
  }

  const result = await getBackendMapAggregates({
    level: query.level,
    parent: query.parent,
    service: query.service,
    q: query.q,
    layers: query.layers,
    viewport,
  });
  if (!result.ok) {
    return NextResponse.json(
      {
        ok: false,
        error: { code: result.code, message: result.message },
        map: emptyMapAggregates(result.message),
      },
      { status: result.status }
    );
  }
  return NextResponse.json({ ok: true, map: result.aggregates });
}

function parseViewport(url: URL) {
  const west = Number(url.searchParams.get("west"));
  const south = Number(url.searchParams.get("south"));
  const east = Number(url.searchParams.get("east"));
  const north = Number(url.searchParams.get("north"));
  if (![west, south, east, north].every(Number.isFinite)) {
    return null;
  }
  return roundViewport({ west, south, east, north });
}
