import type { ReactNode } from "react";

import { PageContainer } from "@/components/layout/page-container";
import { SignedInNav } from "@/features/account/signed-in-nav";
import { OrgSwitcher } from "@/features/business/switcher";
import { loadOrgContext } from "@/lib/business/org";
import { isApiConfigured } from "@/lib/api/backend";
import { readSessionToken } from "@/lib/auth/session";

export default async function BusinessWorkspaceLayout({ children }: { children: ReactNode }) {
  const token = isApiConfigured() ? await readSessionToken() : null;
  const ctx = token ? await loadOrgContext(token) : null;

  return (
    <div className="flex-1">
      <PageContainer className="space-y-3 border-b border-border py-3">
        <SignedInNav kind="business" />
        {ctx ? <OrgSwitcher organizations={ctx.organizations} selected={ctx.organization?.id ?? ""} /> : null}
      </PageContainer>
      {children}
    </div>
  );
}
