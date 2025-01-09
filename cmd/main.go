package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/whatisusername/toon-tank-user-service/api"
	caws "github.com/whatisusername/toon-tank-user-service/internal/aws"
	cconfig "github.com/whatisusername/toon-tank-user-service/internal/config"
	"log/slog"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered", "error", r)
		}
		slog.Info("Program exited")
	}()

	ctx := context.Background()

	secretStore, err := caws.NewSecretsService(ctx)
	if err != nil {
		panic(err)
	}

	cfg, err := cconfig.LoadConfig(ctx, secretStore)
	if err != nil {
		panic(err)
	}

	cognitoAuthSvc, err := caws.NewCognitoService(ctx)
	if err != nil {
		panic(err)
	}

	server, err := api.NewServer(cfg, cognitoAuthSvc)
	if err != nil {
		panic(err)
	}

	lambda.Start(server.HandleRequest)
}
