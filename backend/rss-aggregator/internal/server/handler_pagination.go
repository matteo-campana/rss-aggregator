package server

import (
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultPageSize = 50
	maxPageSize     = 200
)

// pagination is the page window a list endpoint was asked for.
type pagination struct {
	page    int
	perPage int
}

func (p pagination) limit() int32 {
	return int32(p.perPage)
}

func (p pagination) offset() int32 {
	return int32((p.page - 1) * p.perPage)
}

// queryValue reads a trimmed single value out of the query string.
func queryValue(query url.Values, key string) string {
	if values, ok := query[key]; ok && len(values) > 0 {
		return strings.TrimSpace(values[0])
	}

	return ""
}

// parsePagination reads page and per_page, shared by every list endpoint so the
// three of them cannot drift apart. An oversized per_page is clamped rather
// than rejected: the caller still gets a useful answer.
func parsePagination(query url.Values) (pagination, error) {
	parsed := pagination{page: 1, perPage: defaultPageSize}

	if raw := queryValue(query, "page"); raw != "" {
		page, err := strconv.Atoi(raw)
		if err != nil || page < 1 {
			return pagination{}, &badRequestError{"page must be a positive integer"}
		}

		parsed.page = page
	}

	if raw := queryValue(query, "per_page"); raw != "" {
		perPage, err := strconv.Atoi(raw)
		if err != nil || perPage < 1 {
			return pagination{}, &badRequestError{"per_page must be a positive integer"}
		}

		if perPage > maxPageSize {
			perPage = maxPageSize
		}

		parsed.perPage = perPage
	}

	return parsed, nil
}

type badRequestError struct {
	message string
}

func (e *badRequestError) Error() string {
	return e.message
}
