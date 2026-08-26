import { describe, expect, it } from "vitest";

import {
  DEFAULT_PER_PAGE,
  MAX_PER_PAGE,
  parseFilters,
  toApiQuery,
  toBrowserQuery,
  totalPages,
} from "./search-params";

describe("parseFilters", () => {
  it("falls back to the defaults on an empty query string", () => {
    const filters = parseFilters({});

    expect(filters).toEqual({
      search: "",
      category: "",
      minSeeders: null,
      channelId: "",
      sort: "recent",
      page: 1,
      perPage: DEFAULT_PER_PAGE,
    });
  });

  it("reads every filter", () => {
    const filters = parseFilters({
      search: "  SubsPlease  ",
      category: "Anime - English-translated",
      minSeeders: "100",
      channelId: "6a9d1f0e-9e0e-4a4e-9a3b-2d3f4b5c6d7e",
      sort: "seeders",
      page: "3",
      perPage: "25",
    });

    expect(filters.search).toBe("SubsPlease");
    expect(filters.minSeeders).toBe(100);
    expect(filters.sort).toBe("seeders");
    expect(filters.page).toBe(3);
    expect(filters.perPage).toBe(25);
  });

  // A hand-edited URL should show the first page, not an error.
  it.each([
    ["page", "0"],
    ["page", "first"],
    ["perPage", "-1"],
  ])("ignores an invalid %s=%s", (key, value) => {
    const filters = parseFilters({ [key]: value });

    expect(filters.page).toBe(1);
    expect(filters.perPage).toBe(DEFAULT_PER_PAGE);
  });

  it("clamps perPage to the maximum the API accepts", () => {
    expect(parseFilters({ perPage: "100000" }).perPage).toBe(MAX_PER_PAGE);
  });

  it("ignores an unknown sort", () => {
    expect(parseFilters({ sort: "popularity" }).sort).toBe("recent");
  });

  it("ignores a negative minSeeders", () => {
    expect(parseFilters({ minSeeders: "-5" }).minSeeders).toBeNull();
  });

  it("keeps a zero minSeeders, which is a real filter", () => {
    expect(parseFilters({ minSeeders: "0" }).minSeeders).toBe(0);
  });

  it("takes the first value when a parameter is repeated", () => {
    expect(parseFilters({ search: ["one", "two"] }).search).toBe("one");
  });
});

describe("toApiQuery", () => {
  it("renames the parameters to the snake_case the Go API expects", () => {
    const query = toApiQuery(
      parseFilters({ minSeeders: "10", channelId: "abc", perPage: "25", page: "2" }),
    );

    expect(query.get("min_seeders")).toBe("10");
    expect(query.get("channel_id")).toBe("abc");
    expect(query.get("per_page")).toBe("25");
    expect(query.get("page")).toBe("2");
  });

  it("leaves out the filters that are not set", () => {
    const query = toApiQuery(parseFilters({}));

    expect(query.has("search")).toBe(false);
    expect(query.has("min_seeders")).toBe(false);
    expect(query.get("sort")).toBe("recent");
  });
});

describe("toBrowserQuery", () => {
  it("omits the defaults so the URL stays readable", () => {
    expect(toBrowserQuery(parseFilters({}))).toBe("/");
  });

  it("applies the overrides", () => {
    const url = toBrowserQuery(parseFilters({ search: "erai" }), { page: 2 });

    expect(url).toBe("/?search=erai&page=2");
  });
});

describe("totalPages", () => {
  it.each([
    [0, 50, 1],
    [50, 50, 1],
    [51, 50, 2],
    [200, 25, 8],
  ])("total %i over %i per page is %i pages", (total, perPage, expected) => {
    expect(totalPages(total, perPage)).toBe(expected);
  });
});

describe("toBrowserQuery with a base path", () => {
  it("keeps the feeds listing on its own page", () => {
    expect(toBrowserQuery(parseFilters({}), { page: 3 }, "/feeds")).toBe("/feeds?page=3");
  });

  it("returns the bare path when nothing deviates from the defaults", () => {
    expect(toBrowserQuery(parseFilters({}), {}, "/feeds")).toBe("/feeds");
  });

  it("still defaults to the item list", () => {
    expect(toBrowserQuery(parseFilters({}), { page: 2 })).toBe("/?page=2");
  });
});
