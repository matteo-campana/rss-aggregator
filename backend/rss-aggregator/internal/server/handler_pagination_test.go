package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParsePaginationDefaults(t *testing.T) {
	page, err := parsePagination(url.Values{})
	if err != nil {
		t.Fatalf("parsePagination: %v", err)
	}

	if page.page != 1 || page.perPage != defaultPageSize {
		t.Errorf("page/per_page = %d/%d, want 1/%d", page.page, page.perPage, defaultPageSize)
	}

	if page.offset() != 0 || page.limit() != int32(defaultPageSize) {
		t.Errorf("offset/limit = %d/%d, want 0/%d", page.offset(), page.limit(), defaultPageSize)
	}
}

func TestParsePaginationOffset(t *testing.T) {
	page, err := parsePagination(url.Values{"page": {"4"}, "per_page": {"25"}})
	if err != nil {
		t.Fatalf("parsePagination: %v", err)
	}

	// page 4 of 25 starts after the first 75 rows.
	if got, want := page.offset(), int32(75); got != want {
		t.Errorf("offset = %d, want %d", got, want)
	}

	if got, want := page.limit(), int32(25); got != want {
		t.Errorf("limit = %d, want %d", got, want)
	}
}

func TestParsePaginationClampsPageSize(t *testing.T) {
	page, err := parsePagination(url.Values{"per_page": {"9999"}})
	if err != nil {
		t.Fatalf("parsePagination: %v", err)
	}

	if got, want := page.perPage, maxPageSize; got != want {
		t.Errorf("per_page = %d, want it clamped to %d", got, want)
	}
}

func TestParsePaginationRejectsInvalidValues(t *testing.T) {
	for _, query := range []url.Values{
		{"page": {"0"}},
		{"page": {"-1"}},
		{"page": {"first"}},
		{"per_page": {"0"}},
		{"per_page": {"many"}},
	} {
		if _, err := parsePagination(query); err == nil {
			t.Errorf("%v: got nil error, want the request to be rejected", query)
		}
	}
}

// The three list endpoints share parsePagination, so they must all reject the
// same input. Validation runs before the query, so no database is involved.
func TestListEndpointsRejectInvalidPagination(t *testing.T) {
	apiCfg := &ApiConfig{}

	tests := []struct {
		name    string
		handler gin.HandlerFunc
	}{
		{"items", apiCfg.ListItemsHandler()},
		{"feeds", apiCfg.GetFeedsHandler()},
		{"users", apiCfg.GetUsersHandler()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/list", tt.handler)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/list?page=0", nil))

			if got, want := recorder.Code, http.StatusBadRequest; got != want {
				t.Errorf("status = %d, want %d (body: %s)", got, want, recorder.Body.String())
			}
		})
	}
}
