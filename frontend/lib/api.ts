import { ItemFilters, toApiQuery } from "@/lib/search-params";
import { Channel, FeedsPage, ItemsPage, User } from "@/lib/types";

const API_BASE_URL = process.env.API_BASE_URL ?? "http://localhost:8080/api/v1";

/** Thrown when the Go API answers with a non-2xx status. */
export class ApiError extends Error {
  constructor(
    readonly status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function request<T>(path: string, apiKey: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      ...init.headers,
      Authorization: `ApiKey ${apiKey}`,
    },
    // The scraper refreshes the data in the background, so a cached response
    // would show stale seeders.
    cache: "no-store",
  });

  if (!response.ok) {
    const body = await response.text();

    throw new ApiError(response.status, body || response.statusText);
  }

  return (await response.json()) as T;
}

export function getCurrentUser(apiKey: string): Promise<User> {
  return request<User>("/users/me", apiKey);
}

export function listItems(apiKey: string, filters: ItemFilters): Promise<ItemsPage> {
  return request<ItemsPage>(`/items/?${toApiQuery(filters).toString()}`, apiKey);
}

export function listFeeds(apiKey: string, page: number, perPage: number): Promise<FeedsPage> {
  const query = new URLSearchParams({ page: String(page), per_page: String(perPage) });

  return request<FeedsPage>(`/feeds/?${query.toString()}`, apiKey);
}

export function listCategories(apiKey: string): Promise<string[]> {
  return request<{ categories: string[] }>("/items/categories", apiKey).then(
    (body) => body.categories ?? [],
  );
}

export function listChannels(apiKey: string): Promise<Channel[]> {
  return request<Channel[]>("/channels/", apiKey);
}
