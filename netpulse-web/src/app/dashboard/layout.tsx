import type { ReactNode } from "react";

import { SignedInNav } from "@/features/account/signed-in-nav";
import { PageContainer } from "@/components/layout/page-container";

export default function DashboardLayout({ children }: { children: ReactNode }) {
  return (
    <div className="flex-1">
      <PageContainer className="border-b border-border py-3">
        <SignedInNav kind="dashboard" />
      </PageContainer>
      {children}
    </div>
  );
}
