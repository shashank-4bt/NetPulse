import { NextResponse } from "next/server";

import { toShareableReport } from "@/features/intelligence/shareable-report";
import { parseReportId } from "@/lib/reports/id";
import { loadDiagnosis } from "@/lib/reports/load";

export const dynamic = "force-dynamic";

type ReportApiContext = {
  params: Promise<{ id: string }>;
};

export async function GET(request: Request, context: ReportApiContext) {
  const { id } = await context.params;
  const reportId = parseReportId(id);
  if (!reportId) {
    return NextResponse.json(
      { ok: false, error: "Report id is not a UUID." },
      { status: 400 }
    );
  }

  const share = new URL(request.url).searchParams.get("share");
  const loaded = await loadDiagnosis(reportId, share);
  const report = loaded.report;
  if (!report) {
    return NextResponse.json(
      {
        ok: false,
        error:
          "The report is not in this process. That is not a measurement failure.",
      },
      { status: 404 }
    );
  }

  return NextResponse.json({
    ok: true,
    document: toShareableReport(report),
  });
}
