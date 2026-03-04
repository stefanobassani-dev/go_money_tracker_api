package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/stefanobassani-dev/money-tracker/internal/domain"
)

type Service struct {
	jwtService domain.JWTService
	repo       domain.UserRepository
}

func NewAuthService(jwtService domain.JWTService, repo domain.UserRepository) *Service {
	return &Service{
		jwtService: jwtService,
		repo:       repo,
	}
}

func (s *Service) Authenticate(ctx context.Context, email string, password string) (*domain.User, error) {
	user, hash, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrInternal) {
			slog.Error("database error", "err", err)
			return nil, err
		}
		return nil, domain.ErrInvalidAuth
	}

	if !CheckPassword(password, hash) {
		return nil, domain.ErrInvalidAuth
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email string, password string) (*domain.User, string, error) {
	user, err := s.Authenticate(ctx, email, password)
	if err != nil {
		return nil, "", err
	}

	token, err := s.jwtService.Generate(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *Service) Register(ctx context.Context, email string, password string) error {
	_, _, err := s.repo.FindByEmail(ctx, email)

	if err == nil {
		return domain.ErrUserAlreadyExists
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		slog.Error("database error", "err", err)
		return domain.ErrInternal
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return domain.ErrInternal
	}

	return s.repo.CreateUser(ctx, email, hashedPassword)
}
