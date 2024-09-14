package helpers

import (
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
