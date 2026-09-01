import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";

import { createUnavailableReport } from "@/features/diagnose/create-unavailable-report";
import { EvidenceGraph } from "@/features/intelligence/evidence-graph";

describe("EvidenceGraph", () => {
  it("renders every graph node as not measured and opens node details", async () => {
    const user = userEvent.setup();
    const report = createUnavailableReport({
      ok: true,
      raw: "instagram.com",
      hostname: "instagram.com",
      kind: "known_service",
      serviceSlug: "instagram",
    });

    render(<EvidenceGraph report={report} />);

    expect(screen.getByRole("button", { name: /Device/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /Wi-Fi/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /Router/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /Service/i })).toBeInTheDocument();
    expect(screen.getAllByText("Not measured").length).toBeGreaterThanOrEqual(7);
    expect(screen.queryByText("Failed")).not.toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /Device/i }));

    expect(await screen.findByText("No measurement timestamp")).toBeInTheDocument();
    expect(screen.getByText("No numeric value")).toBeInTheDocument();
    expect(screen.getByText("No evidence")).toBeInTheDocument();
    expect(screen.getByText("No measurements")).toBeInTheDocument();
  });
});
