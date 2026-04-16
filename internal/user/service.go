package user

import (
	"API/internal/errs"
	"API/internal/jwt"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Repository interface {
	Registration(ctx context.Context, email string, passwordHash string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}

type ServiceSTR struct {
	repository Repository
	hmacSecret string
}

func NewService(repo Repository, hmacSecret string) *ServiceSTR {
	return &ServiceSTR{repo, hmacSecret}
}

func (s *ServiceSTR) Registration(ctx context.Context, email string, password string) (*Response, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" || password == "" {
		return nil, errs.ErrBadRequest
	}
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user, err := s.repository.Registration(ctx, email, passwordHash)
	if err != nil {
		if errors.Is(err, errs.ErrBadRequest) {
			return nil, errs.ErrBadRequest
		}
		return nil, fmt.Errorf("register user: %w", err)
	}

	return &Response{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
func (s *ServiceSTR) Login(ctx context.Context, email string, password string) (*LoginResponse, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" || password == "" {
		return nil, errs.ErrInvalidCredentials
	}
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	if !checkPasswordHash(password, user.PasswordHash) {
		return nil, errs.ErrInvalidCredentials
	}

	token, err := jwt.Signing(user.ID, s.hmacSecret)
	if err != nil {
		return nil, fmt.Errorf("signing token: %w", err)
	}
	return &LoginResponse{
		User: &Response{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	}, nil
}
func (s *ServiceSTR) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	if id == uuid.Nil {
		return nil, errs.ErrInvalidCredentials
	}
	return s.repository.GetByID(ctx, id)
}
