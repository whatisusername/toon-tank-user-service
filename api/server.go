package api

import (
	"context"
	"errors"
	caws "github.com/whatisusername/toon-tank-user-service/internal/aws"
	cconfig "github.com/whatisusername/toon-tank-user-service/internal/config"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine             *gin.Engine
	ginLambda          *ginadapter.GinLambda
	config             *cconfig.Config
	cognitoAuthService caws.CognitoAuthService
}

func NewServer(cfg *cconfig.Config, cognitoAuthService caws.CognitoAuthService) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}

	s := &Server{
		config:             cfg,
		cognitoAuthService: cognitoAuthService,
	}

	s.registerRoutes()
	return s, nil
}

func (s *Server) registerRoutes() {
	s.engine = gin.Default()

	v1 := s.engine.Group("/v1")
	v1.POST("/users", s.createUser)
	v1.POST("/users/login", s.loginUser)

	s.ginLambda = ginadapter.New(s.engine)

	slog.Info("Routes registered")
}

func (s *Server) HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	slog.Info("Received request", "event", event)
	return s.ginLambda.ProxyWithContext(ctx, event)
}
