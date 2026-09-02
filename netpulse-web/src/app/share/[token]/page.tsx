import type { Metadata } from "next";
import { notFound, redirect } from "next/navigation";

import { getShare } from "@/lib/api/account";
import { isApiConfigured } from "@/lib/api/backend";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Shared report",
  robots: { index: false, follow: false },
};

type PageProps = {
  params: Promise<{ token: string }>;
};

export default async function SharePage({ params }: PageProps) {
  const { token } = await params;
  if (!token || !isApiConfigured()) {
    notFound();
  }
  const loaded = await getShare(token);
  if (!loaded.ok) {
    notFound();
  }
  redirect(`/reports/${loaded.diagnosisId}?share=${encodeURIComponent(token)}`);
}
