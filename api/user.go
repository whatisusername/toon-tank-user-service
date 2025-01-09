package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type createUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if signUpErr := s.cognitoAuthService.SignUp(ctx, s.config.Cognito.ClientID, s.config.Cognito.ClientSecrets, req.Username, req.Password, req.Email); signUpErr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(signUpErr))
		return
	}

	resp := createUserResponse{
		Username: req.Username,
		Email:    req.Email,
	}

	ctx.JSON(http.StatusCreated, successResponse(resp))
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	AccessToken string             `json:"token"`
	User        createUserResponse `json:"user"`
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind request", "error", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cgToken, err := s.cognitoAuthService.Login(ctx, s.config.Cognito.ClientID, s.config.Cognito.ClientSecrets, req.Username, req.Password)
	if err != nil {
		slog.Error("Failed to login", "error", err)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := s.cognitoAuthService.ValidateToken(ctx, s.config.Cognito.UserPoolID, cgToken.AccessToken)
	if err != nil {
		slog.Error("Failed to validate access token", "error", err)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	idToken, err := s.cognitoAuthService.ValidateToken(ctx, s.config.Cognito.UserPoolID, cgToken.IdToken)
	if err != nil {
		slog.Error("Failed to validate id token", "error", err)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	userInfo, err := s.cognitoAuthService.ParseUserInfo(idToken)
	if err != nil {
		slog.Error("Failed to parse user info", "error", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := loginUserResponse{
		AccessToken: accessToken.Raw,
		User: createUserResponse{
			Username: userInfo.Username,
			Email:    userInfo.Email,
		},
	}

	ctx.JSON(http.StatusOK, successResponse(resp))
}
