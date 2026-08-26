// These types mirror the JSON tags of the Go models in
// backend/rss-aggregator/internal/models. Nullable columns are pointers on the
// Go side, so they arrive as null rather than being omitted.

export type Item = {
  id: string;
  title: string | null;
  link: string | null;
  guid: string;
  pubdate: string | null;
  published_at: string | null;
  seeders: number | null;
  leechers: number | null;
  downloads: number | null;
  infohash: string | null;
  category_id: string | null;
  category: string | null;
  size: string | null;
  comments: number | null;
  trusted: string | null;
  remake: string | null;
  description: string | null;
  created_at: string;
  updated_at: string;
  channel_id: string;
  channel_title: string;
};

export type ItemsPage = {
  items: Item[];
  page: number;
  per_page: number;
  total: number;
};

export type Feed = {
  id: string;
  created_at: string;
  updated_at: string;
  url: string;
  name: string;
  last_fetched_at: string | null;
};

export type FeedsPage = {
  feeds: Feed[];
  page: number;
  per_page: number;
  total: number;
};

export type Channel = {
  id: string;
  created_at: string;
  updated_at: string;
  title: string | null;
  description: string | null;
  link: string | null;
  atom_link: string | null;
  feed_id: string | null;
};

export type User = {
  id: string;
  fullname: string;
  firstname: string;
  lastname: string;
  email: string;
};
