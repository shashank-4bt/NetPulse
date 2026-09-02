import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import type { HistoryQuery } from "@/domain/account";

export function HistoryFilters({ query }: { query: HistoryQuery }) {
  return (
    <form method="get" className="grid gap-3 md:grid-cols-5">
      <div>
        <label htmlFor="q" className="text-sm font-medium">
          Search
        </label>
        <Input id="q" name="q" defaultValue={query.q} className="mt-1" />
      </div>
      <div>
        <label htmlFor="target" className="text-sm font-medium">
          Target
        </label>
        <Input id="target" name="target" defaultValue={query.target} className="mt-1" />
      </div>
      <div>
        <label htmlFor="status" className="text-sm font-medium">
          Status
        </label>
        <Input id="status" name="status" defaultValue={query.status} placeholder="queued" className="mt-1" />
      </div>
      <div>
        <label htmlFor="from" className="text-sm font-medium">
          From
        </label>
        <Input id="from" name="from" type="date" defaultValue={query.from} className="mt-1" />
      </div>
      <div>
        <label htmlFor="to" className="text-sm font-medium">
          To
        </label>
        <Input id="to" name="to" type="date" defaultValue={query.to} className="mt-1" />
      </div>
      <div className="md:col-span-5">
        <Button type="submit" variant="outline">
          Filter
        </Button>
      </div>
    </form>
  );
}
