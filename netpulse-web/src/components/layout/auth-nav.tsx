import Link from "next/link";

import { Button } from "@/components/ui/button";
import { getCurrentUser } from "@/lib/auth/session";
import { isApiConfigured } from "@/lib/api/backend";

export async function AuthNav() {
  const user = isApiConfigured() ? await getCurrentUser() : null;
  if (user) {
    return (
      <Button nativeButton={false} variant="outline" render={<Link href="/dashboard" />}>
        Dashboard
      </Button>
    );
  }
  return (
    <Button nativeButton={false} variant="outline" render={<Link href="/login" />}>
      Sign in
    </Button>
  );
}
