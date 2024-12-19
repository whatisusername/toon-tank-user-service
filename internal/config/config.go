package config

import (
	"context"
	"encoding/json"

	caws "github.com/whatisusername/toon-tank-user-service/internal/aws"
	"github.com/whatisusername/toon-tank-user-service/internal/env"
)

type CognitoConfig struct {
	UserPoolID    string `json:"userPoolId"`
	ClientID      string `json:"clientId"`
	ClientSecrets string `json:"clientSecrets"`
}

type Config struct {
	Cognito CognitoConfig `json:"cognito"`
}

func LoadConfig(ctx context.Context, secretStore caws.SecretStore) (*Config, error) {
	secretName, err := env.GetValue("SECRET_NAME")
	if err != nil {
		return nil, err
	}

	secrets, err := secretStore.GetSecretValue(ctx, secretName)
	if err != nil {
		return nil, err
	}

	var cc CognitoConfig
	if err = json.Unmarshal([]byte(*secrets), &cc); err != nil {
		return nil, err
	}

	return &Config{
		Cognito: cc,
	}, nil
}
