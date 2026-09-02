import type { HistoryQuery } from "@/domain/account";

export function parseHistoryQuery(
  params: Record<string, string | string[] | undefined>
): HistoryQuery {
  return {
    q: first(params.q),
    status: first(params.status),
    target: first(params.target),
    from: first(params.from),
    to: first(params.to),
  };
}

export function historySearchParams(query: HistoryQuery): string {
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(query)) {
    if (value) {
      params.set(key, value);
    }
  }
  return params.toString();
}

function first(value: string | string[] | undefined): string {
  if (Array.isArray(value)) {
    return value[0] ?? "";
  }
  return value?.trim() ?? "";
}
