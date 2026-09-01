"use client";

import { useState } from "react";
import { ChevronRightIcon } from "lucide-react";

import { EmptyState } from "@/components/feedback/empty-state";
import { ConfidenceBadge } from "@/components/status/confidence-badge";
import { LayerStatusBadge } from "@/components/status/layer-status-badge";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import type { DiagnosticReport, EvidenceGraphNode } from "@/domain/diagnostic";
import { formatConfidenceValue } from "@/features/intelligence/confidence";
import { lookupByIds } from "@/features/intelligence/honesty";

type EvidenceGraphProps = {
  report: DiagnosticReport;
};

export function EvidenceGraph({ report }: EvidenceGraphProps) {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const selected = report.graph.find((node) => node.id === selectedId) ?? null;

  return (
    <section aria-labelledby="graph-heading" className="space-y-4">
      <div className="space-y-1">
        <h2 id="graph-heading" className="text-lg font-semibold">
          Evidence graph
        </h2>
        <p className="text-sm text-muted-foreground">
          Device → Wi-Fi → Router → ISP → Route → CDN → Service. Open a node to
          inspect evidence, measurements, status, confidence, and timestamp.
        </p>
      </div>
      <ol className="flex flex-col gap-2 md:flex-row md:flex-wrap md:items-stretch">
        {report.graph.map((node, index) => (
          <li key={node.id} className="flex min-w-0 flex-1 flex-col gap-2 md:flex-row md:items-stretch">
            <Button
              type="button"
              variant="outline"
              onClick={() => setSelectedId(node.id)}
              aria-haspopup="dialog"
              className="h-auto min-h-11 w-full flex-col items-start gap-2 whitespace-normal px-3 py-3 text-left"
            >
              <span className="font-mono text-xs text-muted-foreground">
                {String(index + 1).padStart(2, "0")}
              </span>
              <span className="text-sm font-medium text-foreground">
                {node.label}
              </span>
              <LayerStatusBadge status={node.status} />
            </Button>
            {index < report.graph.length - 1 ? (
              <span
                className="flex items-center justify-center px-1 text-muted-foreground"
                aria-hidden="true"
              >
                <ChevronRightIcon className="hidden size-4 md:block" />
                <span className="md:hidden">↓</span>
              </span>
            ) : null}
          </li>
        ))}
      </ol>
      <NodeDialog
        node={selected}
        report={report}
        open={selected !== null}
        onOpenChange={(open) => {
          if (!open) {
            setSelectedId(null);
          }
        }}
      />
    </section>
  );
}

type NodeDialogProps = {
  node: EvidenceGraphNode | null;
  report: DiagnosticReport;
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

function NodeDialog({ node, report, open, onOpenChange }: NodeDialogProps) {
  const evidence = node
    ? lookupByIds(report.evidence, node.evidenceIds)
    : [];
  const measurements = node
    ? lookupByIds(report.measurements, node.measurementIds)
    : [];

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg" showCloseButton>
        {node ? (
          <>
            <DialogHeader>
              <DialogTitle>{node.label}</DialogTitle>
              <DialogDescription>
                Evidence and measurements for this layer. Empty fields mean the
                probe did not run.
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4">
              <div className="flex flex-wrap items-center gap-2">
                <LayerStatusBadge status={node.status} />
                {node.confidence.level ? (
                  <ConfidenceBadge confidence={node.confidence.level} />
                ) : (
                  <Badge variant="outline">Confidence not supplied</Badge>
                )}
              </div>
              <dl className="grid gap-3 text-sm">
                <div>
                  <dt className="font-medium">Confidence value</dt>
                  <dd className="mt-1 font-mono text-muted-foreground">
                    {formatConfidenceValue(node.confidence)}
                  </dd>
                </div>
                <div>
                  <dt className="font-medium">Timestamp</dt>
                  <dd className="mt-1 font-mono text-muted-foreground">
                    {node.timestamp ?? "No measurement timestamp"}
                  </dd>
                </div>
              </dl>
              <div>
                <h3 className="text-sm font-medium">Evidence</h3>
                {evidence.length === 0 ? (
                  <EmptyState
                    className="mt-2"
                    title="No evidence"
                    description="This node has no attached measured facts."
                  />
                ) : (
                  <ul className="mt-2 space-y-2 text-sm text-muted-foreground">
                    {evidence.map((item) => (
                      <li key={item.id}>
                        <span className="font-medium text-foreground">
                          {item.title}.
                        </span>{" "}
                        {item.body}
                      </li>
                    ))}
                  </ul>
                )}
              </div>
              <div>
                <h3 className="text-sm font-medium">Measurements</h3>
                {measurements.length === 0 ? (
                  <EmptyState
                    className="mt-2"
                    title="No measurements"
                    description="This node has no attached measurement records."
                  />
                ) : (
                  <ul className="mt-2 space-y-1 font-mono text-xs text-muted-foreground">
                    {measurements.map((item) => (
                      <li key={item.id}>
                        {item.label}:{" "}
                        {item.measured
                          ? `${item.value ?? "—"} ${item.unit ?? ""}`.trim()
                          : "not measured"}
                      </li>
                    ))}
                  </ul>
                )}
              </div>
            </div>
          </>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}
