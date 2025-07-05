package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type ServiceAccount struct {
	Type           string `json:"type"`
	ProjectId      string `json:"project_id"`
	PrivateKeyId   string `json:"private_key_id"`
	PrivateKey     string `json:"private_key"`
	ClientEmail    string `json:"client_email"`
	ClientId       string `json:"client_id"`
	AuthUri        string `json:"auth_uri"`
	TokenUri       string `json:"token_uri"`
	AuthProvider   string `json:"auth_provider_x509_cert_url"`
	ClientCertUrl  string `json:"client_x509_cert_url"`
	UniverseDomain string `json:"universe_domain"`
}

var Config ServiceAccount

func LoadConfig() error {
	env := os.Getenv("NODE_ENV")
	svcAccount := os.Getenv("GCP_SERVICE_ACCOUNT")

	if env == "" || svcAccount == "" {
		return fmt.Errorf("env or service account not provided")
	}

	decodedData, err := base64.StdEncoding.DecodeString(svcAccount)
	if err != nil {
		return fmt.Errorf("error load base64: %v", err.Error())
	}

	cfg := &ServiceAccount{}
	if err := json.Unmarshal(decodedData, cfg); err != nil {
		return fmt.Errorf("failed to unmarshall : %v", err.Error())
	}

	Config = *cfg

	return nil
}
