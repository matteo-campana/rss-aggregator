package auth

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		want    string
		wantErr error
	}{
		{name: "valid", header: "ApiKey secret-key", want: "secret-key"},
		{name: "extra spaces", header: "ApiKey    secret-key", want: "secret-key"},
		{name: "missing", header: "", wantErr: ErrNoAuthHeader},
		{name: "no scheme", header: "secret-key", wantErr: ErrMalformedAuthHeader},
		{name: "wrong scheme", header: "Bearer secret-key", wantErr: ErrMalformedAuthHeader},
		{name: "no key", header: "ApiKey", wantErr: ErrMalformedAuthHeader},
		{name: "too many fields", header: "ApiKey secret-key extra", wantErr: ErrMalformedAuthHeader},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.header != "" {
				headers.Set("Authorization", tt.header)
			}

			got, err := GetAPIKey(headers)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("key = %q, want %q", got, tt.want)
			}
		})
	}
}
