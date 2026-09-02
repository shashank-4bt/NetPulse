const SAFE_METHODS = new Set(["GET", "HEAD", "OPTIONS"]);

export function csrfAllowed(request: Request): boolean {
  if (SAFE_METHODS.has(request.method.toUpperCase())) {
    return true;
  }
  const expected = new URL(request.url).origin;
  const origin = request.headers.get("origin");
  if (origin) {
    return sameOrigin(origin, expected);
  }
  const referer = request.headers.get("referer");
  if (referer) {
    return sameOrigin(referer, expected);
  }
  return request.headers.get("x-netpulse-csrf") === "same-origin";
}

function sameOrigin(value: string, expected: string): boolean {
  try {
    return new URL(value).origin === expected;
  } catch {
    return false;
  }
}
