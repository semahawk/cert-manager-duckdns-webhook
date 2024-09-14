package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	"github.com/pkg/errors"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"os"
	"strings"

	"net/http"

	"k8s.io/client-go/rest"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
)

var GroupName = os.Getenv("GROUP_NAME")

var httpClient = &http.Client{}

const (
	duckDNSSuffix       = "duckdns.org"
	missingSecretKeyMsg = "key %q not found in secret \"%s/%s\""
	missingSecretErrMsg = "Unable to load secret %v %q"
	duckDNSUpdateURL    = "https://www.duckdns.org/update"
)

func main() {
	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	cmd.RunWebhookServer(GroupName,
		&duckDNSProviderSolver{},
	)
}

type duckDNSProviderSolver struct {
	client *kubernetes.Clientset
}

type duckDNSProviderConfig struct {
	APITokenSecretRef cmmeta.SecretKeySelector `json:"apiTokenSecretRef"`
}

func (c *duckDNSProviderSolver) Name() string {
	return "duckdns"
}

func UpdateTXTRecord(ctx context.Context, domain string, token string, txt string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, duckDNSUpdateURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("domains", domain)
	q.Add("token", token)
	q.Add("dns_txt", txt)
	req.URL.RawQuery = q.Encode()

	return httpClient.Do(req)
}

func getDomainName(DNSName string) string {
	var result string

	domainName := strings.TrimSuffix(DNSName, duckDNSSuffix) // prefix.domain. or domain.
	// Remove trailing dot
	domainName = util.UnFqdn(domainName) // prefix.domain or domain

	split := strings.Split(domainName, ".")

	// If it's prefix.domain, return domain
	if len(split) == 2 {
		result = split[1]
	} else {
		result = domainName
	}
	klog.Infof("Got domain %v from DNS %v", result, DNSName)

	return result
}

func (c *duckDNSProviderSolver) getDuckDNSToken(cfg *duckDNSProviderConfig, namespace string) (*string, error) {
	secretName := cfg.APITokenSecretRef.LocalObjectReference.Name

	secret, err := c.client.CoreV1().Secrets(namespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, missingSecretErrMsg, secretName, namespace+"/"+secretName)
	}

	data, ok := secret.Data[cfg.APITokenSecretRef.Key]
	if !ok {
		return nil, fmt.Errorf(missingSecretKeyMsg, cfg.APITokenSecretRef.Key,
			cfg.APITokenSecretRef.LocalObjectReference.Name, namespace)
	}

	apiKey := string(data)
	return &apiKey, nil
}

func getToken(solver duckDNSProviderSolver, ch *v1alpha1.ChallengeRequest) (*string, error) {
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return nil, err
	}

	token, err := solver.getDuckDNSToken(&cfg, ch.ResourceNamespace)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (c *duckDNSProviderSolver) fetchTokenAndDomain(ch *v1alpha1.ChallengeRequest) (string, string, error) {
	token, err := getToken(*c, ch)
	if err != nil {
		klog.Errorf("Unable to get token: %v", err)
		return "", "", err
	}
	domain := getDomainName(ch.DNSName)
	return domain, *token, nil
}

func (c *duckDNSProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	domain, token, err := c.fetchTokenAndDomain(ch)
	if err != nil {
		return err
	}
	res, err := UpdateTXTRecord(context.Background(), domain, token, ch.Key)
	if err != nil {
		klog.Errorf("Failed to set TXT record: %v", err)
		return err
	}
	klog.Infof("Got status %v when updating TXT record for %v", res.Status, ch.ResolvedFQDN)
	return nil
}

func (c *duckDNSProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	domain, token, err := c.fetchTokenAndDomain(ch)
	if err != nil {
		return err
	}
	_, err = UpdateTXTRecord(context.Background(), domain, token, "")
	if err != nil {
		klog.Errorf("Failed to clean TXT record for DNS %v: %v", ch.DNSName, err)
		return err
	}
	return nil
}

func (c *duckDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, _ <-chan struct{}) error {
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}
	c.client = cl

	return nil
}

func loadConfig(cfgJSON *extapi.JSON) (duckDNSProviderConfig, error) {
	cfg := duckDNSProviderConfig{}

	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}

	return cfg, nil
}

func (c *duckDNSProviderSolver) validateConfig(cfg *duckDNSProviderConfig) error {

	if cfg.APITokenSecretRef.LocalObjectReference.Name == "" {
		return errors.New("No api token secret provided in DuckDNS config")
	}

	return nil
}
