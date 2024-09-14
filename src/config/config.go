package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"k8s.io/klog/v2"

	"github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type DuckDNSProviderConfig struct {
	APITokenSecretRef v1.SecretKeySelector `json:"apiTokenSecretRef"`
}

func LoadConfig(cfgJSON *extapi.JSON) (DuckDNSProviderConfig, error) {
	cfg := DuckDNSProviderConfig{}

	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}

	if err := validateConfig(&cfg); err != nil {
		klog.Errorf("Invalid config: %v", err)
		return cfg, err
	}

	return cfg, nil
}

func validateConfig(cfg *DuckDNSProviderConfig) error {
	if cfg.APITokenSecretRef.LocalObjectReference.Name == "" {
		return errors.New("no api token secret provided in DuckDNS config")
	}

	return nil
}
