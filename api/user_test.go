package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	caws "github.com/whatisusername/toon-tank-user-service/internal/aws"
	cconfig "github.com/whatisusername/toon-tank-user-service/internal/config"
)

func TestServer_createUser(t *testing.T) {
	tests := []struct {
		name          string
		body          gin.H
		buildStubs    func(authSvc *caws.MockCognitoAuthService)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": "test",
				"email":    "test@example.com",
				"password": "test123456A",
			},
			buildStubs: func(authSvc *caws.MockCognitoAuthService) {
				authSvc.EXPECT().
					SignUp(mock.Anything, "fake_client_id", "fake_client_secret", "test", "test123456A", "test@example.com").
					Return(nil).Once()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				expected, err := json.Marshal(response{
					Success: true,
					Message: "Success",
					Data: createUserResponse{
						Username: "test",
						Email:    "test@example.com",
					},
				})
				assert.NoError(t, err)
				assert.Equal(t, expected, recorder.Body.Bytes())
			},
		},
		{
			name: "Bad Request",
			body: gin.H{
				"username": "test",
				"password": "test123456A",
			},
			buildStubs: func(authSvc *caws.MockCognitoAuthService) {
				authSvc.AssertNotCalled(t, "SignUp")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"username": "test",
				"email":    "test@example.com",
				"password": "test123456A",
			},
			buildStubs: func(authSvc *caws.MockCognitoAuthService) {
				authSvc.EXPECT().
					SignUp(mock.Anything, "fake_client_id", "fake_client_secret", "test", "test123456A", "test@example.com").
					Return(fmt.Errorf("server is busy")).Once()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cognitoAuthService := caws.NewMockCognitoAuthService(t)
			tt.buildStubs(cognitoAuthService)

			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			url := "/v1/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			testServer := newTestServer(t, cognitoAuthService)
			recorder := httptest.NewRecorder()

			testServer.engine.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func TestServer_loginUser(t *testing.T) {
	fakeToken := &caws.CognitoToken{
		IdToken:      "fake_id_token",
		AccessToken:  "fake_access_token",
		RefreshToken: "fake_refresh_token",
	}

	tests := []struct {
		name          string
		body          gin.H
		buildStubs    func(authSvc *caws.MockCognitoAuthService)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": "test",
				"password": "test123456A",
			},
			buildStubs: func(authSvc *caws.MockCognitoAuthService) {
				authSvc.EXPECT().
					Login(mock.Anything, "fake_client_id", "fake_client_secret", "test", "test123456A").
					Return(fakeToken, nil).Once()

				mockTokenValidation(authSvc, "us-east-1_example", "fake_access_token", jwt.MapClaims{"email": "test@example.com"})
				mockTokenValidation(authSvc, "us-east-1_example", "fake_id_token", jwt.MapClaims{"cognito:username": "test", "email": "test@example.com"})

				authSvc.EXPECT().ParseUserInfo(mock.AnythingOfType("*jwt.Token")).
					Return(&caws.CognitoUserInfo{
						Username: "test",
						Email:    "test@example.com",
					}, nil).Once()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				expected, err := json.Marshal(response{
					Success: true,
					Message: "Success",
					Data: loginUserResponse{
						AccessToken: "fake_access_token",
						User: createUserResponse{
							Username: "test",
							Email:    "test@example.com",
						},
					},
				})
				assert.NoError(t, err)
				assert.Equal(t, expected, recorder.Body.Bytes())
			},
		},
		{
			name: "Bad Request",
			body: gin.H{
				"username": "test",
			},
			buildStubs: func(authSvc *caws.MockCognitoAuthService) {
				authSvc.AssertNotCalled(t, "Login")
				authSvc.AssertNotCalled(t, "ValidateToken")
				authSvc.AssertNotCalled(t, "ParseUserInfo")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Login Failed",
			body: gin.H{
				"username": "test",
				"password": "test123456",
			},
			buildStubs: func(authSvc *caws.MockCognitoAuthService) {
				authSvc.EXPECT().
					Login(mock.Anything, "fake_client_id", "fake_client_secret", "test", "test123456").
					Return(nil, errors.New("incorrect username or password")).Once()

				authSvc.AssertNotCalled(t, "ValidateToken")
				authSvc.AssertNotCalled(t, "ParseUserInfo")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid Access Token",
			body: gin.H{
				"username": "test",
				"password": "test123456A",
			},
			buildStubs: func(authSvc *caws.MockCognitoAuthService) {
				authSvc.EXPECT().
					Login(mock.Anything, "fake_client_id", "fake_client_secret", "test", "test123456A").
					Return(fakeToken, nil).Once()

				authSvc.EXPECT().ValidateToken(mock.Anything, "us-east-1_example", "fake_access_token").
					Return(nil, errors.New("invalid access token")).Once()

				authSvc.AssertNotCalled(t, "ParseUserInfo")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid Id Token",
			body: gin.H{
				"username": "test",
				"password": "test123456A",
			},
			buildStubs: func(authSvc *caws.MockCognitoAuthService) {
				authSvc.EXPECT().
					Login(mock.Anything, "fake_client_id", "fake_client_secret", "test", "test123456A").
					Return(fakeToken, nil).Once()

				mockTokenValidation(authSvc, "us-east-1_example", "fake_access_token", jwt.MapClaims{"email": "test@example.com"})

				authSvc.EXPECT().ValidateToken(mock.Anything, "us-east-1_example", "fake_id_token").
					Return(nil, errors.New("invalid id token")).Once()

				authSvc.AssertNotCalled(t, "ParseUserInfo")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Parse User Info Failed",
			body: gin.H{
				"username": "test",
				"password": "test123456A",
			},
			buildStubs: func(authSvc *caws.MockCognitoAuthService) {
				authSvc.EXPECT().
					Login(mock.Anything, "fake_client_id", "fake_client_secret", "test", "test123456A").
					Return(fakeToken, nil).Once()

				mockTokenValidation(authSvc, "us-east-1_example", "fake_access_token", jwt.MapClaims{"email": "test@example.com"})
				mockTokenValidation(authSvc, "us-east-1_example", "fake_id_token", jwt.MapClaims{"cognito:username": "test", "email": "test@example.com"})

				authSvc.EXPECT().ParseUserInfo(mock.AnythingOfType("*jwt.Token")).
					Return(nil, errors.New("unexpected claims type")).Once()
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cognitoAuthService := caws.NewMockCognitoAuthService(t)
			tt.buildStubs(cognitoAuthService)

			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			url := "/v1/users/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			testServer := newTestServer(t, cognitoAuthService)
			recorder := httptest.NewRecorder()

			testServer.engine.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func newTestServer(t *testing.T, cognitoAuthService *caws.MockCognitoAuthService) *Server {
	t.Helper()
	cfg := &cconfig.Config{
		Cognito: cconfig.CognitoConfig{
			UserPoolID:    "us-east-1_example",
			ClientID:      "fake_client_id",
			ClientSecrets: "fake_client_secret",
		}}

	server, err := NewServer(cfg, cognitoAuthService)
	require.NoError(t, err)
	return server
}

func mockTokenValidation(authSvc *caws.MockCognitoAuthService, userPoolId, token string, claims jwt.MapClaims) {
	authSvc.EXPECT().ValidateToken(mock.Anything, userPoolId, token).
		Return(&jwt.Token{
			Raw:       token,
			Method:    &jwt.SigningMethodRSA{},
			Header:    map[string]interface{}{"alg": "RS256", "kid": "fake_kid"},
			Claims:    jwt.Claims(claims),
			Signature: []byte{},
			Valid:     true,
		}, nil).Once()
}
