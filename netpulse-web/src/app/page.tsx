import type { Metadata } from "next";

import { DiagnoseSection } from "@/components/home/diagnose-section";
import { HealthSection } from "@/components/home/health-section";
import { HeroSection } from "@/components/home/hero-section";
import { IncidentsSection } from "@/components/home/incidents-section";
import {
  AudiencePreview,
  EvidencePreview,
  FinalCtaSection,
  HowItWorksPreview,
  IntelligencePreview,
  TrustPreview,
} from "@/components/home/method-sections";
import { ServicesSection } from "@/components/home/services-section";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Is the internet down?",
  description:
    "NetPulse finds out what's actually broken — device, network, ISP, or service — using measured evidence.",
  alternates: { canonical: "/" },
};

export default function HomePage() {
  return (
    <main id="main-content" className="flex-1">
      <HeroSection />
      <DiagnoseSection />
      <HealthSection />
      <IncidentsSection />
      <ServicesSection />
      <HowItWorksPreview />
      <EvidencePreview />
      <IntelligencePreview />
      <AudiencePreview />
      <TrustPreview />
      <FinalCtaSection />
    </main>
  );
}
