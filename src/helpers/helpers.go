package helpers

import (
	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	"io"
	"k8s.io/klog/v2"
	"net/http"
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

func GetResponseBody(res *http.Response) (string, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		klog.Errorf("Unable to get body from response")
		return "", err
	}
	result := string(body)
	result = strings.ReplaceAll(result, "\n", "\t")
	return result, err
}
