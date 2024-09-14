package helpers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDomainName(t *testing.T) {
	tests := []struct {
		name         string
		DNSName      string
		expectedName string
	}{
		{
			name:         "Full domain with subdomain",
			DNSName:      "subdomain.example.duckdns.org",
			expectedName: "example",
		},
		{
			name:         "Domain without subdomain",
			DNSName:      "example.duckdns.org",
			expectedName: "example",
		},
		{
			name:         "Edge case with no suffix",
			DNSName:      "example",
			expectedName: "example",
		},
		{
			name:         "Empty domain",
			DNSName:      "",
			expectedName: "",
		},
		{
			name:         "Domain with trailing dot",
			DNSName:      "example.duckdns.org.",
			expectedName: "example",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDomainName(tt.DNSName)
			assert.Equal(t, tt.expectedName, result, "Expected domain name to match for test case: %s", tt.name)
		})
	}
}

func TestGetResponseBody(t *testing.T) {
	tests := []struct {
		name         string
		bodyReader   io.Reader
		expectedBody string
		expectError  bool
	}{
		{
			name:         "Valid response body",
			bodyReader:   strings.NewReader("test body"),
			expectedBody: "test body",
			expectError:  false,
		},
		{
			name:         "Body with new line",
			bodyReader:   strings.NewReader("test\nbody"),
			expectedBody: "test\tbody",
			expectError:  false,
		},
		{
			name:         "Empty response body",
			bodyReader:   strings.NewReader(""),
			expectedBody: "",
			expectError:  false,
		},
		{
			name:         "Error reading response body",
			bodyReader:   nil,
			expectedBody: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{
				Body: io.NopCloser(tt.bodyReader),
			}

			if tt.bodyReader == nil {
				res.Body = &errorReader{}
			}

			body, err := GetResponseBody(res)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, body)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("simulated read error")
}

func (e *errorReader) Close() error {
	return nil
}
