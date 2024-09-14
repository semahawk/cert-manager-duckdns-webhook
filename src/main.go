package main

import (
	"context"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/cert-manager/webhook-example/src/config"
	"github.com/cert-manager/webhook-example/src/duckdns"
	"github.com/cert-manager/webhook-example/src/helpers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"os"
)

var GroupName = os.Getenv("GROUP_NAME")

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

func (solver *duckDNSProviderSolver) Name() string {
	return "duckdns"
}

func (solver *duckDNSProviderSolver) createDuckDNSConnector(ch *v1alpha1.ChallengeRequest) (*duckdns.Connector, string, error) {
	cfg, err := config.LoadConfig(ch.Config)
	if err != nil {
		klog.Errorf("Unable to load config: %v", err)
		return nil, "", err
	}

	token, err := config.GetDuckDNSToken(solver.client, &cfg, ch.ResourceNamespace)
	if err != nil {
		klog.Errorf("Unable to get token: %v", err)
		return nil, "", err
	}

	domain := helpers.GetDomainName(ch.DNSName)

	connector := duckdns.NewConnector(*token)
	return connector, domain, nil
}

func (solver *duckDNSProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	connector, domain, err := solver.createDuckDNSConnector(ch)
	if err != nil {
		return err
	}

	res, err := connector.SetTXTRecord(context.Background(), domain, ch.Key)
	if err != nil {
		klog.Errorf("Failed to set TXT record: %v", err)
		return err
	}
	klog.Infof("Got status %v when updating TXT record for %v", res.Status, ch.ResolvedFQDN)
	return nil
}

func (solver *duckDNSProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	connector, domain, err := solver.createDuckDNSConnector(ch)
	if err != nil {
		return err
	}

	_, err = connector.CleanTXTRecord(context.Background(), domain)
	if err != nil {
		klog.Errorf("Failed to clean TXT record for DNS %v: %v", ch.DNSName, err)
		return err
	}
	return nil
}

func (solver *duckDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, _ <-chan struct{}) error {
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}
	solver.client = cl

	return nil
}
