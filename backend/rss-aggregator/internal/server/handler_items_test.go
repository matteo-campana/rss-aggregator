package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParseListItemsQueryDefaults(t *testing.T) {
	query, err := parseListItemsQuery(url.Values{})
	if err != nil {
		t.Fatalf("parseListItemsQuery: %v", err)
	}

	if got, want := query.page, 1; got != want {
		t.Errorf("page = %d, want %d", got, want)
	}

	if got, want := query.perPage, defaultPageSize; got != want {
		t.Errorf("per_page = %d, want %d", got, want)
	}

	if got, want := query.params.Sort, sortRecent; got != want {
		t.Errorf("sort = %q, want %q", got, want)
	}

	if got, want := query.params.PageOffset, int32(0); got != want {
		t.Errorf("offset = %d, want %d", got, want)
	}

	for name, valid := range map[string]bool{
		"search":      query.params.Search.Valid,
		"category":    query.params.Category.Valid,
		"min_seeders": query.params.MinSeeders.Valid,
		"channel_id":  query.params.ChannelID.Valid,
	} {
		if valid {
			t.Errorf("%s is set, want it unset so the filter is skipped", name)
		}
	}
}

func TestParseListItemsQueryFilters(t *testing.T) {
	query, err := parseListItemsQuery(url.Values{
		"search":      {"  SubsPlease  "},
		"category":    {"Anime - English-translated"},
		"min_seeders": {"100"},
		"channel_id":  {"6a9d1f0e-9e0e-4a4e-9a3b-2d3f4b5c6d7e"},
		"sort":        {"seeders"},
		"page":        {"3"},
		"per_page":    {"25"},
	})

	if err != nil {
		t.Fatalf("parseListItemsQuery: %v", err)
	}

	if got, want := query.params.Search.String, "SubsPlease"; got != want {
		t.Errorf("search = %q, want %q (surrounding spaces trimmed)", got, want)
	}

	if got, want := query.params.MinSeeders.Int32, int32(100); got != want {
		t.Errorf("min_seeders = %d, want %d", got, want)
	}

	if !query.params.ChannelID.Valid {
		t.Error("channel_id was not parsed")
	}

	// page 3 of 25 starts after the first 50 rows.
	if got, want := query.params.PageOffset, int32(50); got != want {
		t.Errorf("offset = %d, want %d", got, want)
	}

	if got, want := query.params.PageSize, int32(25); got != want {
		t.Errorf("page size = %d, want %d", got, want)
	}
}

func TestParseListItemsQueryClampsPageSize(t *testing.T) {
	query, err := parseListItemsQuery(url.Values{"per_page": {"100000"}})
	if err != nil {
		t.Fatalf("parseListItemsQuery: %v", err)
	}

	if got, want := query.perPage, maxPageSize; got != want {
		t.Errorf("per_page = %d, want it clamped to %d", got, want)
	}
}

func TestParseListItemsQueryRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		query url.Values
	}{
		{name: "page zero", query: url.Values{"page": {"0"}}},
		{name: "page not a number", query: url.Values{"page": {"first"}}},
		{name: "per_page zero", query: url.Values{"per_page": {"0"}}},
		{name: "negative min_seeders", query: url.Values{"min_seeders": {"-1"}}},
		{name: "min_seeders not a number", query: url.Values{"min_seeders": {"many"}}},
		{name: "unknown sort", query: url.Values{"sort": {"popularity"}}},
		{name: "malformed channel id", query: url.Values{"channel_id": {"not-a-uuid"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseListItemsQuery(tt.query); err == nil {
				t.Fatal("got nil error, want the request to be rejected")
			}
		})
	}
}

// The validation runs before the query, so a bad request needs no database.
func TestListItemsHandlerRejectsInvalidQuery(t *testing.T) {
	apiCfg := &ApiConfig{}

	router := gin.New()
	router.GET("/items", apiCfg.ListItemsHandler())

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/items?sort=popularity", nil))

	if got, want := recorder.Code, http.StatusBadRequest; got != want {
		t.Errorf("status = %d, want %d (body: %s)", got, want, recorder.Body.String())
	}
}
