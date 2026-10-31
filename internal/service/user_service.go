package service

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/storage"
	"advanced-blog-management-system/pkg/auth"
	"context"
	"strings"
)

type UserService struct {
	userStore  *storage.UserStore
	jwtManager *auth.JWTManager
}

func NewUserService(userStore *storage.UserStore, jwtManager *auth.JWTManager) *UserService {
	return &UserService{userStore: userStore, jwtManager: jwtManager}
}

func (s *UserService) Register(ctx context.Context, req *model.UserCreateRequest) (*model.TokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	username := strings.TrimSpace(req.Username)

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u, err := s.userStore.Create(username, email, hash)
	if err != nil {
		return nil, apperrors.ErrUserAlreadyExists
	}

	token, expiresAt, err := s.jwtManager.GenerateToken(u.ID, u.Email, u.Username)
	if err != nil {
		return nil, err
	}
	return &model.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: model.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		},
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.TokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	u, err := s.userStore.GetByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}
	if !auth.CheckPassword(req.Password, u.Password) {
		return nil, apperrors.ErrInvalidCredentials
	}
	token, expiresAt, err := s.jwtManager.GenerateToken(u.ID, u.Email, u.Username)
	if err != nil {
		return nil, err
	}
	return &model.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: model.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		},
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	u, err := s.userStore.GetByID(id)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}
	return &model.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	u, err := s.userStore.GetByEmail(email)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}
	return &model.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}
