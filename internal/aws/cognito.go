package aws

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/whatisusername/toon-tank-user-service/internal/token"
)

type CognitoToken struct {
	IdToken      string
	AccessToken  string
	RefreshToken string
}

type CognitoUserInfo struct {
	Username string
	Email    string
}

type CognitoAuthService interface {
	SignUp(ctx context.Context, clientId, clientSecret, username, password, email string) error
	Login(ctx context.Context, clientId, clientSecret, username, password string) (*CognitoToken, error)
	ValidateToken(ctx context.Context, userPoolId, tokenString string) (*jwt.Token, error)
	ParseUserInfo(idToken *jwt.Token) (*CognitoUserInfo, error)
}

func NewCognitoService(ctx context.Context, optFns ...func(options *config.LoadOptions) error) (*CognitoService, error) {
	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, err
	}

	client := cognitoidentityprovider.NewFromConfig(cfg)

	return &CognitoService{
		client: client,
	}, nil
}

type CognitoService struct {
	client *cognitoidentityprovider.Client
}

func (c *CognitoService) SignUp(ctx context.Context, clientId, clientSecret, username, password, email string) error {
	slog.Info("Signing up user", "username", username)

	secretHash, err := token.GenerateBase64HMAC(clientSecret, username+clientId)
	if err != nil {
		return err
	}

	_, err = c.client.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(clientId),
		Username:   aws.String(username),
		Password:   aws.String(password),
		SecretHash: aws.String(secretHash),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	})

	if err != nil {
		return err
	}

	slog.Info("Created user", "username", username)

	return nil
}

func (c *CognitoService) Login(ctx context.Context, clientId, clientSecret, username, password string) (*CognitoToken, error) {
	slog.Info("Logging in user", "username", username)

	secretHash, err := token.GenerateBase64HMAC(clientSecret, username+clientId)
	if err != nil {
		return nil, err
	}

	output, err := c.client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(clientId),
		AuthParameters: map[string]string{
			"USERNAME":    username,
			"PASSWORD":    password,
			"SECRET_HASH": secretHash,
		},
	})

	if err != nil {
		return nil, err
	}

	slog.Info("Logged in user", "token type", *output.AuthenticationResult.TokenType, "expires in", output.AuthenticationResult.ExpiresIn)

	return &CognitoToken{
		IdToken:      *output.AuthenticationResult.IdToken,
		AccessToken:  *output.AuthenticationResult.AccessToken,
		RefreshToken: *output.AuthenticationResult.RefreshToken,
	}, nil
}

func (c *CognitoService) ValidateToken(ctx context.Context, userPoolId, tokenString string) (*jwt.Token, error) {
	region := strings.Split(userPoolId, "_")[0]
	jwkUrl := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolId)
	k, err := keyfunc.NewDefault([]string{jwkUrl})
	if err != nil {
		return nil, fmt.Errorf("failed to load JWK: %w", err)
	}

	t, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return k.Keyfunc(token)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return t, nil
}

func (c *CognitoService) ParseUserInfo(idToken *jwt.Token) (*CognitoUserInfo, error) {
	var (
		claims   jwt.MapClaims
		ok       bool
		raw      interface{}
		username string
		email    string
	)

	claims, ok = idToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unexpected claims type: %T", idToken.Claims)
	}

	if raw, ok = claims["cognito:username"]; ok {
		if username, ok = raw.(string); !ok {
			return nil, fmt.Errorf("username claim is not a string")
		}
	}

	if raw, ok = claims["email"]; ok {
		if email, ok = raw.(string); !ok {
			return nil, fmt.Errorf("email claim is not a string")
		}
	}

	return &CognitoUserInfo{
		Username: username,
		Email:    email,
	}, nil
}
