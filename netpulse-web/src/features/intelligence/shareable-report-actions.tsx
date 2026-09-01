"use client";

import { useState } from "react";

import { Button } from "@/components/ui/button";
import type { DiagnosticReport } from "@/domain/diagnostic";
import { stringifyShareableReport } from "@/features/intelligence/shareable-report";

type ShareableReportActionsProps = {
  report: DiagnosticReport;
};

export function ShareableReportActions({ report }: ShareableReportActionsProps) {
  const [status, setStatus] = useState<"idle" | "copied" | "failed">("idle");

  async function copyJson() {
    setStatus("idle");
    try {
      await navigator.clipboard.writeText(stringifyShareableReport(report));
      setStatus("copied");
    } catch {
      setStatus("failed");
    }
  }

  return (
    <div className="flex flex-wrap items-center gap-3">
      <Button type="button" variant="outline" onClick={copyJson}>
        Copy report JSON
      </Button>
      <Button type="button" variant="outline" onClick={() => window.print()}>
        Print / save as PDF
      </Button>
      {status === "copied" ? (
        <p className="text-sm text-muted-foreground" role="status">
          Copied the shareable report document.
        </p>
      ) : null}
      {status === "failed" ? (
        <p className="text-sm text-destructive" role="alert">
          Clipboard is unavailable. The JSON was not copied.
        </p>
      ) : null}
    </div>
  );
}
