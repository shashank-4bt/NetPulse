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
    default: `${publicConfig.appName} — ${publicConfig.appTagline}`,
    template: `%s · ${publicConfig.appName}`,
  },
  description:
    "NetPulse isolates whether an internet problem is your device, network, ISP, or the service — using measured evidence, not guesses.",
  robots: {
    index: true,
    follow: true,
  },
  openGraph: {
    title: `${publicConfig.appName} — ${publicConfig.appTagline}`,
    description:
      "Find out what's actually broken. Evidence-based internet diagnosis.",
    type: "website",
    siteName: publicConfig.appName,
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
