export const API_ERROR_CODES = [
  "validation_error",
  "unauthorized",
  "forbidden",
  "not_found",
  "rate_limited",
  "unavailable",
  "ssrf_blocked",
  "internal",
  "timeout",
  "unknown",
] as const;

export type ApiErrorCode = (typeof API_ERROR_CODES)[number];

export type ApiFailure = {
  ok: false;
  code: ApiErrorCode;
  message: string;
  status: number;
};

export function apiFailure(
  code: ApiErrorCode,
  message: string,
  status: number
): ApiFailure {
  return { ok: false, code, message, status };
}

export function classifyTransportError(error: unknown): ApiFailure {
  if (error instanceof DOMException && error.name === "AbortError") {
    return apiFailure(
      "timeout",
      "The diagnose service timed out. No substitute result was generated.",
      504
    );
  }
  return apiFailure(
    "unavailable",
    "The diagnose service is unavailable. No fallback diagnosis was invented.",
    503
  );
}

export function parseErrorCode(value: unknown): ApiErrorCode {
  if (typeof value === "string" && (API_ERROR_CODES as readonly string[]).includes(value)) {
    return value as ApiErrorCode;
  }
  return "unknown";
}
