import Link from "next/link";

import { PageContainer } from "@/components/layout/page-container";
import { Button } from "@/components/ui/button";

export default function NotFound() {
  return (
    <main id="main-content" className="flex flex-1 items-center">
      <PageContainer className="py-16">
        <h1 className="text-2xl font-semibold tracking-tight">Page not found</h1>
        <p className="mt-2 max-w-md text-muted-foreground">
          That URL is not part of the public site.
        </p>
        <Button className="mt-6" nativeButton={false} render={<Link href="/" />}>
          Back to home
        </Button>
      </PageContainer>
    </main>
  );
}
