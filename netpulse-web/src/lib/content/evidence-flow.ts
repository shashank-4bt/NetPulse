export const EVIDENCE_FLOW = [
  {
    id: "evidence",
    name: "Evidence",
    summary: "Recorded observations from probes. These are measured facts.",
  },
  {
    id: "hypothesis",
    name: "Hypothesis",
    summary:
      "An inference that cites evidence. Never shown as a measured fact.",
  },
  {
    id: "confidence",
    name: "Confidence",
    summary: "How strongly the evidence supports the hypothesis.",
  },
  {
    id: "recommendation",
    name: "Recommendation",
    summary: "A safe next step. Not automatically executed.",
  },
] as const;
