import { FlaskConicalIcon } from "lucide-react";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

type DevelopmentBannerProps = {
  title?: string;
  description?: string;
};

export function DevelopmentBanner({
  title = "Development state",
  description = "Live internet health, outages, and maps are unavailable. Nothing on this page is a live measurement.",
}: DevelopmentBannerProps) {
  return (
    <Alert className="border-border bg-muted/40">
      <FlaskConicalIcon />
      <AlertTitle>{title}</AlertTitle>
      <AlertDescription>{description}</AlertDescription>
    </Alert>
  );
}
