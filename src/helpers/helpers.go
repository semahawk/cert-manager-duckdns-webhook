package helpers

import (
	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	"strings"
)

const duckDNSSuffix = ".duckdns.org"

func GetDomainName(DNSName string) string {
	domainName := util.UnFqdn(DNSName) // Remove trailing dot (domain.duckdns.org. -> domain.duckdns.org)
	domainName = strings.TrimSuffix(domainName, duckDNSSuffix)

	split := strings.Split(domainName, ".")

	// If it's prefix.domain, return domain
	if len(split) == 2 {
		return split[1]
	} else {
		return domainName
	}
}
