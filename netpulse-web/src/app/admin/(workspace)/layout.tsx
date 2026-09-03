import { SignedInNav } from "@/features/account/signed-in-nav";
import { PageContainer } from "@/components/layout/page-container";
import type { ReactNode } from "react";

export default function AdminWorkspaceLayout({ children }: { children: ReactNode }) {
  return (
    <div className="flex-1">
      <PageContainer className="space-y-3 border-b border-border py-3">
        <SignedInNav kind="admin" />
      </PageContainer>
      {children}
    </div>
  );
}
