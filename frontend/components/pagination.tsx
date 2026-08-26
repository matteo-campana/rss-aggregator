import Link from "next/link";

import { ItemFilters, toBrowserQuery, totalPages } from "@/lib/search-params";

type Props = {
  filters: ItemFilters;
  total: number;
};

export function Pagination({ filters, total }: Props) {
  const pages = totalPages(total, filters.perPage);

  if (pages <= 1) {
    return null;
  }

  const previous = filters.page > 1 ? toBrowserQuery(filters, { page: filters.page - 1 }) : null;
  const next = filters.page < pages ? toBrowserQuery(filters, { page: filters.page + 1 }) : null;

  const linkClass =
    "rounded border border-neutral-300 px-3 py-1.5 text-sm dark:border-neutral-700";

  return (
    <nav className="flex items-center gap-3" aria-label="Pagination">
      {previous ? (
        <Link href={previous} className={linkClass}>
          Previous
        </Link>
      ) : (
        <span className={`${linkClass} opacity-40`}>Previous</span>
      )}

      <span className="text-sm text-neutral-600 dark:text-neutral-400">
        Page {filters.page} of {pages}
      </span>

      {next ? (
        <Link href={next} className={linkClass}>
          Next
        </Link>
      ) : (
        <span className={`${linkClass} opacity-40`}>Next</span>
      )}
    </nav>
  );
}
