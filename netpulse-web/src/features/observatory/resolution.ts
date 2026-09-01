export type ResolutionInput = {
  recoverySampleCount: number;
  identifiedCause: boolean;
};

export type ResolutionDecision = {
  ok: boolean;
  reason: string;
};

export function canMarkResolved(input: ResolutionInput): ResolutionDecision {
  if (input.recoverySampleCount < 2) {
    return {
      ok: false,
      reason:
        "One recovered measurement is not enough to mark an incident resolved.",
    };
  }
  if (!input.identifiedCause) {
    return {
      ok: false,
      reason:
        "Resolution requires an identified cause backed by evidence, not recovery alone.",
    };
  }
  return {
    ok: true,
    reason:
      "Independent recoveries and an identified cause are present. Resolution is still a stored judgment, not a population claim.",
  };
}
