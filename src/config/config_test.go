package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestLoadConfig_ValidConfig(t *testing.T) {
	cfgJSON := &extapi.JSON{
		Raw: json.RawMessage(`{
			"apiTokenSecretRef": {
				"name": "my-secret",
				"key": "token"
			}
		}`),
	}

	cfg, err := LoadConfig(cfgJSON)

	assert.NoError(t, err, "Expected no error for valid config")
	assert.Equal(t, "my-secret", cfg.APITokenSecretRef.LocalObjectReference.Name, "Expected API token secret name to match")
	assert.Equal(t, "token", cfg.APITokenSecretRef.Key, "Expected API token secret key to match")
}

func TestLoadConfig_InvalidConfigName(t *testing.T) {
	cfgJSON := &extapi.JSON{
		Raw: json.RawMessage(`{
			"apiTokenSecretRef": {
				"name": "",
				"key": "token"
			}
		}`),
	}

	_, err := LoadConfig(cfgJSON)
	assert.Error(t, err, "Expected error for invalid config")
	assert.Contains(t, err.Error(), "no api token secret name provided in DuckDNS config", "Expected specific error message for missing API token secret")
}

func TestLoadConfig_InvalidConfiKey(t *testing.T) {
	cfgJSON := &extapi.JSON{
		Raw: json.RawMessage(`{
			"apiTokenSecretRef": {
				"name": "my-secret",
				"key": ""
			}
		}`),
	}

	_, err := LoadConfig(cfgJSON)
	assert.Error(t, err, "Expected error for invalid config")
	assert.Contains(t, err.Error(), "no api token secret key provided in DuckDNS config", "Expected specific error message for missing API token secret")
}
