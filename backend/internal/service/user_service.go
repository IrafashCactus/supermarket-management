package service

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/supermarket/backend/internal/auth"
	"github.com/supermarket/backend/internal/database/sqlc"
	"github.com/supermarket/backend/internal/dto"
)

type UserService struct {
	queries sqlc.Querier
}

func NewUserService(queries sqlc.Querier) *UserService {
	return &UserService{queries: queries}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         req.Role,
	})
	if err != nil {
		return dto.UserResponse{}, err
	}

	return mapToUserResponse(user), nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (dto.UserResponse, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return mapToUserResponse(user), nil
}

func (s *UserService) ListUsers(ctx context.Context, page, perPage int) (dto.ListResponse[dto.UserResponse], error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	users, err := s.queries.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(perPage),
		Offset: int32(offset),
	})
	if err != nil {
		return dto.ListResponse[dto.UserResponse]{}, err
	}

	total, err := s.queries.CountUsers(ctx)
	if err != nil {
		return dto.ListResponse[dto.UserResponse]{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	var items []dto.UserResponse
	for _, u := range users {
		items = append(items, mapToUserResponse(u))
	}

	return dto.ListResponse[dto.UserResponse]{
		Data: items,
		Meta: dto.PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	user, err := s.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		IsActive: req.IsActive,
	})
	if err != nil {
		return dto.UserResponse{}, err
	}

	return mapToUserResponse(user), nil
}

func (s *UserService) ResetPassword(ctx context.Context, id uuid.UUID, req dto.ResetPasswordRequest) error {
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	_, err = s.queries.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: hash,
	})
	return err
}

func mapToUserResponse(user sqlc.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
