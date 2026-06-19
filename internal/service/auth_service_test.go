package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	db "bitedash/internal/db/sqlc"
	"bitedash/internal/pkg/auth"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAuthQuerier struct {
	createUser     func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error)
	getUserByEmail func(ctx context.Context, email string) (db.User, error)
	getUserByID    func(ctx context.Context, id uuid.UUID) (db.User, error)
}

func (m *mockAuthQuerier) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
	if m.createUser != nil {
		return m.createUser(ctx, arg)
	}
	return db.CreateUserRow{}, nil
}

func (m *mockAuthQuerier) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	if m.getUserByEmail != nil {
		return m.getUserByEmail(ctx, email)
	}
	return db.User{}, nil
}

func (m *mockAuthQuerier) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	if m.getUserByID != nil {
		return m.getUserByID(ctx, id)
	}
	return db.User{}, nil
}

func TestAuthService_Register_success(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	mock := &mockAuthQuerier{
		createUser: func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
			assert.Equal(t, "user@example.com", arg.Email)
			assert.NotEmpty(t, arg.PasswordHash)
			return db.CreateUserRow{ID: userID, Email: arg.Email, CreatedAt: time.Now()}, nil
		},
	}

	svc := NewAuthService(mock)
	token, err := svc.Register(context.Background(), "user@example.com", "password123")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_Register_duplicateEmail(t *testing.T) {
	t.Parallel()

	mock := &mockAuthQuerier{
		createUser: func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
			return db.CreateUserRow{}, &pgconn.PgError{Code: "23505"}
		},
	}

	svc := NewAuthService(mock)
	_, err := svc.Register(context.Background(), "dup@example.com", "password123")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrEmailAlreadyExists)
}

func TestAuthService_Register_dbError(t *testing.T) {
	t.Parallel()

	dbErr := errors.New("connection lost")
	mock := &mockAuthQuerier{
		createUser: func(ctx context.Context, arg db.CreateUserParams) (db.CreateUserRow, error) {
			return db.CreateUserRow{}, dbErr
		},
	}

	svc := NewAuthService(mock)
	_, err := svc.Register(context.Background(), "user@example.com", "password123")
	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
}

func TestAuthService_Login_success(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	hash, err := auth.HashPassword("secret")
	require.NoError(t, err)

	mock := &mockAuthQuerier{
		getUserByEmail: func(ctx context.Context, email string) (db.User, error) {
			return db.User{
				ID:           userID,
				Email:        email,
				PasswordHash: hash,
				CreatedAt:    time.Now(),
			}, nil
		},
	}

	svc := NewAuthService(mock)
	token, err := svc.Login(context.Background(), "user@example.com", "secret")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_Login_userNotFound(t *testing.T) {
	t.Parallel()

	mock := &mockAuthQuerier{
		getUserByEmail: func(ctx context.Context, email string) (db.User, error) {
			return db.User{}, sql.ErrNoRows
		},
	}

	svc := NewAuthService(mock)
	_, err := svc.Login(context.Background(), "missing@example.com", "secret")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Login_wrongPassword(t *testing.T) {
	t.Parallel()

	hash, err := auth.HashPassword("correct")
	require.NoError(t, err)

	mock := &mockAuthQuerier{
		getUserByEmail: func(ctx context.Context, email string) (db.User, error) {
			return db.User{ID: uuid.New(), Email: email, PasswordHash: hash}, nil
		},
	}

	svc := NewAuthService(mock)
	_, err = svc.Login(context.Background(), "user@example.com", "wrong")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}
