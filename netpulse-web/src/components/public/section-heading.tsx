import { cn } from "@/lib/utils";

type SectionHeadingProps = {
  id: string;
  eyebrow?: string;
  title: string;
  description?: string;
  className?: string;
};

export function SectionHeading({
  id,
  eyebrow,
  title,
  description,
  className,
}: SectionHeadingProps) {
  return (
    <div className={cn("max-w-2xl", className)}>
      {eyebrow ? (
        <p className="text-xs font-medium tracking-wide text-muted-foreground uppercase">
          {eyebrow}
        </p>
      ) : null}
      <h2 id={id} className="mt-1 text-xl font-semibold tracking-tight md:text-2xl">
        {title}
      </h2>
      {description ? (
        <p className="mt-2 text-sm text-muted-foreground md:text-base">
          {description}
        </p>
      ) : null}
    </div>
  );
}
