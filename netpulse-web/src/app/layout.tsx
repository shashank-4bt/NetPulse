import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";

import { AppProviders } from "@/components/layout/app-providers";
import { SiteFooter } from "@/components/layout/site-footer";
import { SiteHeader } from "@/components/layout/site-header";
import { SkipLink } from "@/components/layout/skip-link";
import { publicConfig } from "@/lib/config/public";

import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  metadataBase: new URL(
    process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000"
  ),
  title: {
    default: `${publicConfig.appName} · Design system`,
    template: `%s · ${publicConfig.appName}`,
  },
  description:
    "NetPulse design system foundation for the Internet Health Intelligence Platform. Product diagnosis is not available yet.",
  robots: {
    index: false,
    follow: false,
  },
  openGraph: {
    title: `${publicConfig.appName} · Design system`,
    description:
      "Foundational UI for NetPulse. No live internet health data is shown.",
    type: "website",
  },
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className={`${geistSans.variable} ${geistMono.variable} h-full`}
    >
      <body className="flex min-h-full flex-col">
        <AppProviders>
          <SkipLink />
          <SiteHeader />
          <div className="flex flex-1 flex-col">{children}</div>
          <SiteFooter />
        </AppProviders>
      </body>
    </html>
  );
}
