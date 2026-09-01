import { DiagnoseForm } from "@/features/diagnose/diagnose-form";
import { PageContainer } from "@/components/layout/page-container";
import { SectionContainer } from "@/components/layout/section-container";
import { SectionHeading } from "@/components/public/section-heading";

export function DiagnoseSection() {
  return (
    <SectionContainer id="diagnose" labelledBy="diagnose-heading">
      <PageContainer>
        <SectionHeading
          id="diagnose-heading"
          eyebrow="Quick diagnose"
          title="Start a diagnosis"
          description="The full workflow lives on /diagnose. Submitting here validates the target and opens that run. Probe steps stay unavailable until workers exist."
        />
        <div className="mt-6 rounded-lg border border-border bg-card p-4 md:p-6">
          <DiagnoseForm />
        </div>
      </PageContainer>
    </SectionContainer>
  );
}
