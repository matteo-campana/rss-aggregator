export function EmptyState({ filtered }: { filtered: boolean }) {
  return (
    <div className="rounded border border-dashed border-neutral-300 p-10 text-center dark:border-neutral-700">
      <p className="text-sm font-medium">No items to show.</p>

      <p className="mx-auto mt-2 max-w-prose text-sm text-neutral-600 dark:text-neutral-400">
        {filtered
          ? "No stored item matches these filters. Try widening them or resetting."
          : "The database has no items yet. Add a feed, then run the API with SCRAPER_ENABLED=true so the background scraper can collect them."}
      </p>
    </div>
  );
}
