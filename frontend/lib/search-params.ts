// Filter state lives in the URL, so a filtered view can be linked and shared.
// Everything here is pure, which keeps it testable without a browser.

export const SORTS = ["recent", "seeders", "oldest"] as const;

export type Sort = (typeof SORTS)[number];

export const DEFAULT_PER_PAGE = 50;
export const MAX_PER_PAGE = 200;

export type ItemFilters = {
  search: string;
  category: string;
  minSeeders: number | null;
  channelId: string;
  sort: Sort;
  page: number;
  perPage: number;
};

export type RawSearchParams = Record<string, string | string[] | undefined>;

function first(value: string | string[] | undefined): string {
  if (Array.isArray(value)) {
    return value[0] ?? "";
  }

  return value ?? "";
}

function positiveInt(value: string, fallback: number): number {
  const parsed = Number.parseInt(value, 10);

  if (!Number.isFinite(parsed) || parsed < 1) {
    return fallback;
  }

  return parsed;
}

/**
 * Reads the query string into filters, falling back to the defaults on
 * anything malformed: a bad URL should show the first page, not an error.
 */
export function parseFilters(params: RawSearchParams): ItemFilters {
  const sort = first(params.sort) as Sort;
  const rawMinSeeders = Number.parseInt(first(params.minSeeders), 10);

  return {
    search: first(params.search).trim(),
    category: first(params.category).trim(),
    minSeeders:
      Number.isFinite(rawMinSeeders) && rawMinSeeders >= 0 ? rawMinSeeders : null,
    channelId: first(params.channelId).trim(),
    sort: SORTS.includes(sort) ? sort : "recent",
    page: positiveInt(first(params.page), 1),
    perPage: Math.min(positiveInt(first(params.perPage), DEFAULT_PER_PAGE), MAX_PER_PAGE),
  };
}

/** Builds the query string the Go API expects. */
export function toApiQuery(filters: ItemFilters): URLSearchParams {
  const query = new URLSearchParams({
    sort: filters.sort,
    page: String(filters.page),
    per_page: String(filters.perPage),
  });

  if (filters.search) {
    query.set("search", filters.search);
  }

  if (filters.category) {
    query.set("category", filters.category);
  }

  if (filters.minSeeders !== null) {
    query.set("min_seeders", String(filters.minSeeders));
  }

  if (filters.channelId) {
    query.set("channel_id", filters.channelId);
  }

  return query;
}

/**
 * Builds the query string for a link within the app. Defaults are left out so
 * the URL stays short and readable.
 */
export function toBrowserQuery(filters: ItemFilters, overrides: Partial<ItemFilters> = {}): string {
  const merged = { ...filters, ...overrides };
  const query = new URLSearchParams();

  if (merged.search) {
    query.set("search", merged.search);
  }

  if (merged.category) {
    query.set("category", merged.category);
  }

  if (merged.minSeeders !== null) {
    query.set("minSeeders", String(merged.minSeeders));
  }

  if (merged.channelId) {
    query.set("channelId", merged.channelId);
  }

  if (merged.sort !== "recent") {
    query.set("sort", merged.sort);
  }

  if (merged.page > 1) {
    query.set("page", String(merged.page));
  }

  if (merged.perPage !== DEFAULT_PER_PAGE) {
    query.set("perPage", String(merged.perPage));
  }

  const serialised = query.toString();

  return serialised ? `/?${serialised}` : "/";
}

export function totalPages(total: number, perPage: number): number {
  return Math.max(1, Math.ceil(total / perPage));
}
