package config

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	caws "github.com/whatisusername/toon-tank-user-service/internal/aws"
)

func TestLoadConfig(t *testing.T) {
	cognitoCfg := fakeCognitoConfig()
	cognitoCfgString := fakeCognitoConfigString(t, cognitoCfg)

	tests := []struct {
		name                    string
		setupEnv                func(t *testing.T)
		mockSecretStoreResponse func(secretStore *caws.MockSecretStore)
		want                    *Config
		wantErr                 bool
	}{
		{
			name: "Valid Config",
			setupEnv: func(t *testing.T) {
				t.Setenv("SECRET_NAME", "test")
			},
			mockSecretStoreResponse: func(secretStore *caws.MockSecretStore) {
				secretStore.EXPECT().
					GetSecretValue(mock.Anything, mock.AnythingOfType("string")).
					Return(&cognitoCfgString, nil).
					Once()
			},
			want:    &Config{Cognito: *cognitoCfg},
			wantErr: false,
		},
		{
			name:     "Missing Env Variable",
			setupEnv: func(t *testing.T) {},
			mockSecretStoreResponse: func(secretStore *caws.MockSecretStore) {
				secretStore.AssertNotCalled(t, "GetSecretValue")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid Secrets String",
			setupEnv: func(t *testing.T) {
				t.Setenv("SECRET_NAME", "test")
			},
			mockSecretStoreResponse: func(secretStore *caws.MockSecretStore) {
				secretStore.EXPECT().
					GetSecretValue(mock.Anything, mock.AnythingOfType("string")).
					Return(stringPtr("\"userPoolId\"=\"us-east-1_example\""), nil).
					Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Secrets Manager Unreachable",
			setupEnv: func(t *testing.T) {
				t.Setenv("SECRET_NAME", "test")
			},
			mockSecretStoreResponse: func(secretStore *caws.MockSecretStore) {
				secretStore.EXPECT().
					GetSecretValue(mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("timeout")).
					Once()
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv(t)

			secretStore := caws.NewMockSecretStore(t)
			tt.mockSecretStoreResponse(secretStore)

			got, err := LoadConfig(context.Background(), secretStore)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error: %v", err)
			}
			assert.Equal(t, tt.want, got, "LoadConfig() returned unexpected result")
		})
	}
}

func fakeCognitoConfig() *CognitoConfig {
	return &CognitoConfig{
		UserPoolID:    "us-east-1_example",
		ClientID:      "fake_client_id",
		ClientSecrets: "fake_client_secret",
	}
}

func fakeCognitoConfigString(t *testing.T, cognito *CognitoConfig) string {
	data, err := json.Marshal(cognito)
	require.NoError(t, err)

	return string(data)
}

func stringPtr(s string) *string {
	return &s
}
