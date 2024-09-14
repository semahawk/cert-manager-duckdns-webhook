package helpers

import (
	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	"strings"
)

const duckDNSSuffix = "duckdns.org"

func GetDomainName(DNSName string) string {
	domainName := strings.TrimSuffix(DNSName, duckDNSSuffix) // prefix.domain. or domain.
	// Remove trailing dot
	domainName = util.UnFqdn(domainName) // prefix.domain or domain

	split := strings.Split(domainName, ".")

	// If it's prefix.domain, return domain
	if len(split) == 2 {
		return split[1]
	} else {
		return domainName
	}
}
