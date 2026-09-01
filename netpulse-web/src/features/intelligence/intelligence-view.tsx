import { ConfidencePanel } from "@/features/intelligence/confidence-panel";
import { EscalationList } from "@/features/intelligence/escalation-list";
import { EvidenceGraph } from "@/features/intelligence/evidence-graph";
import { RecommendationList } from "@/features/intelligence/recommendation-list";
import { ShareableReportActions } from "@/features/intelligence/shareable-report-actions";
import type { DiagnosticReport } from "@/domain/diagnostic";

type IntelligenceViewProps = {
  report: DiagnosticReport;
};

export function IntelligenceView({ report }: IntelligenceViewProps) {
  return (
    <div className="space-y-10">
      <EvidenceGraph report={report} />
      <ConfidencePanel
        confidence={report.confidence}
        evidence={report.evidence}
        alternatives={report.alternativeHypotheses}
      />
      <RecommendationList recommendations={report.recommendations} />
      <EscalationList conditions={report.escalationConditions} />
      <ShareableReportActions report={report} />
    </div>
  );
}
