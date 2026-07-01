package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "bitedash/internal/db/sqlc"
	"bitedash/internal/pkg/auth"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"go.opentelemetry.io/otel"
)

var authTracer = otel.Tracer("bitedash/internal/service/auth")

type AuthService struct {
	queries authQuerier
}

func NewAuthService(q authQuerier) *AuthService {
	return &AuthService{queries: q}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (string, error) {
	ctx, span := authTracer.Start(ctx, "AuthService.Register")
	defer span.End()

	hash, err := auth.HashPassword(password)
	if err != nil {
		return "", err
	}

	userID := uuid.New()

	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           userID,
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	})

	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			if pgErr.Code == "23505" {
				return "", ErrEmailAlreadyExists
			}

		}

		return "", err
	}

	token, err := auth.GenerateToken(user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	ctx, span := authTracer.Start(ctx, "AuthService.Login")
	defer span.End()

	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = auth.CheckPasswordHash(password, user.PasswordHash)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	ctx, span := authTracer.Start(ctx, "AuthService.GetUserByID")
	defer span.End()

	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, ErrUserNotFound
		}
		return db.User{}, err
	}
	return user, nil
}
