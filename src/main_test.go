package main

import (
	"github.com/csp33/cert-manager-duckdns-webhook/src/duckdns"
	"os"
	"testing"

	acmetest "github.com/cert-manager/cert-manager/test/acme"
)

var (
	domain = os.Getenv("DUCKDNS_DOMAIN_URL")
	zone   = domain + "."
)

func TestRunsSuite(t *testing.T) {
	solver := duckdns.NewSolver(nil)
	fixture := acmetest.NewFixture(solver,
		acmetest.SetResolvedZone(zone),
		acmetest.SetDNSName(domain),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetManifestPath("../testdata/duckdns-solver"),
	)

	fixture.RunBasic(t)
}
