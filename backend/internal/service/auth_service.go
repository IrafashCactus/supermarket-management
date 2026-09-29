package service

import (
	"context"
	"errors"

	"github.com/supermarket/backend/internal/auth"
	"github.com/supermarket/backend/internal/database/sqlc"
	"github.com/supermarket/backend/internal/dto"
	"github.com/google/uuid"
)

type AuthService struct {
	queries    sqlc.Querier
	jwtManager *auth.JWTManager
}

func NewAuthService(queries sqlc.Querier, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		queries:    queries,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (dto.LoginResponse, string, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return dto.LoginResponse{}, "", errors.New("invalid email or password")
	}

	if !user.IsActive {
		return dto.LoginResponse{}, "", errors.New("account is disabled")
	}

	if err := auth.CheckPassword(password, user.PasswordHash); err != nil {
		return dto.LoginResponse{}, "", errors.New("invalid email or password")
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return dto.LoginResponse{}, "", err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return dto.LoginResponse{}, "", err
	}

	resp := dto.LoginResponse{
		AccessToken: accessToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return resp, refreshToken, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (string, error) {
	claims, err := s.jwtManager.ValidateToken(refreshTokenString)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return "", errors.New("invalid token subject")
	}

	user, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !user.IsActive {
		return "", errors.New("account is disabled")
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
