import { NextResponse } from "next/server";

import { startDiagnosis } from "@/features/diagnose/start-diagnosis";
import { toShareableReport } from "@/features/intelligence/shareable-report";
import { loadDiagnosis } from "@/lib/reports/load";

export const dynamic = "force-dynamic";

export async function POST(request: Request) {
  let body: unknown;
  try {
    body = await request.json();
  } catch {
    return NextResponse.json(
      { ok: false, error: "Expected JSON with a target." },
      { status: 400 }
    );
  }

  const target =
    typeof body === "object" && body !== null && "target" in body
      ? String((body as { target: unknown }).target)
      : "";

  const result = await startDiagnosis(target);
  if (!result.ok) {
    return NextResponse.json(
      { ok: false, error: result.error, outcome: "invalid_input" },
      { status: 400 }
    );
  }

  const report = (await loadDiagnosis(result.reportId)).report;
  if (!report) {
    return NextResponse.json(
      {
        ok: false,
        error: "The report was not stored in this process.",
        outcome: "backend_unavailable",
      },
      { status: 503 }
    );
  }

  return NextResponse.json(
    { ok: true, document: toShareableReport(report) },
    { status: 201 }
  );
}
