package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/localstack"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretsService_GetSecretValue(t *testing.T) {
	ctx := context.Background()

	container, err := localstack.Run(ctx, "localstack/localstack:4.0.3")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx), "failed to terminate container")
	})

	host, err := container.Host(ctx)
	require.NoError(t, err)

	mappedPort, err := container.MappedPort(ctx, "4566/tcp")
	require.NoError(t, err)

	endpoint := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())

	svc, err := NewSecretsService(ctx, func(o *secretsmanager.Options) {
		o.Region = "us-east-1"
		o.Credentials = aws.AnonymousCredentials{}
		o.BaseEndpoint = &endpoint
	})
	require.NoError(t, err)

	_, err = svc.client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String("test-secret"),
		SecretString: aws.String("test-secret-value"),
	})
	require.NoError(t, err)

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				key: "test-secret",
			},
			want:    aws.String("test-secret-value"),
			wantErr: false,
		},
		{
			name: "Not Found",
			args: args{
				key: "test-secret-1",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.GetSecretValue(ctx, tt.args.key)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "unexpected error: %v", err)
			}

			if tt.want == nil {
				assert.Nil(t, got, "GetSecretValue() returned unexpected result")
			} else {
				assert.Equal(t, *tt.want, *got, "GetSecretValue() returned unexpected result")
			}
		})
	}
}
