import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { StatusBadge } from "@/components/status/status-badge";
import { OPERATIONAL_STATUSES } from "@/domain/display";
import { STATUS_LABEL } from "@/lib/design/taxonomy";

describe("StatusBadge", () => {
  it("renders a text label for every operational status", () => {
    for (const status of OPERATIONAL_STATUSES) {
      const { unmount } = render(<StatusBadge status={status} />);
      expect(screen.getByText(STATUS_LABEL[status])).toBeInTheDocument();
      unmount();
    }
  });
});
