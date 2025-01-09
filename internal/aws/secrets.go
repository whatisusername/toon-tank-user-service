package aws

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretStore interface {
	GetSecretValue(ctx context.Context, key string) (*string, error)
}

type SecretsService struct {
	client *secretsmanager.Client
}

func NewSecretsService(ctx context.Context, optFns ...func(*secretsmanager.Options)) (*SecretsService, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg, optFns...)

	return &SecretsService{
		client: client,
	}, nil
}

func (s *SecretsService) GetSecretValue(ctx context.Context, key string) (*string, error) {
	output, err := s.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	slog.Info("Got secret", "secret", key, "metadata", output.ResultMetadata)

	return output.SecretString, nil
}
