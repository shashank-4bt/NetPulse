const MAX_PREFILL = 2048;

export function safeDiagnosePrefill(value: string | undefined): string {
  if (!value) {
    return "";
  }
  const trimmed = value.trim();
  if (!trimmed || trimmed.length > MAX_PREFILL) {
    return "";
  }
  if (/[<>'"\\]/.test(trimmed)) {
    return "";
  }
  return trimmed;
}
