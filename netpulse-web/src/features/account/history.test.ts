import { describe, expect, it } from "vitest";

import { historySearchParams, parseHistoryQuery } from "@/features/account/history";

describe("history query", () => {
  it("reads filters without inventing a date range", () => {
    expect(parseHistoryQuery({ q: "youtube", status: "queued" })).toEqual({
      q: "youtube",
      status: "queued",
      target: "",
      from: "",
      to: "",
    });
  });

  it("omits empty filters from the query string", () => {
    expect(
      historySearchParams({
        q: "dns",
        status: "",
        target: "example.com",
        from: "",
        to: "",
      })
    ).toBe("q=dns&target=example.com");
  });
});
