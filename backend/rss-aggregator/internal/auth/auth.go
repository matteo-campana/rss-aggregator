// Package auth reads the credentials carried by an incoming request.
package auth

import (
	"errors"
	"net/http"
	"strings"
)

// ErrNoAuthHeader is returned when the request carries no Authorization header.
var ErrNoAuthHeader = errors.New("authorization header is missing")

// ErrMalformedAuthHeader is returned when the header is present but does not
// follow the expected `ApiKey <key>` form.
var ErrMalformedAuthHeader = errors.New("malformed authorization header, expected: ApiKey <key>")

// GetAPIKey extracts the API key from the request headers.
// Example:
//
//	Authorization: ApiKey 123456
func GetAPIKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", ErrNoAuthHeader
	}

	fields := strings.Fields(authorization)
	if len(fields) != 2 || fields[0] != "ApiKey" || fields[1] == "" {
		return "", ErrMalformedAuthHeader
	}

	return fields[1], nil
}
