package config

import "errors"

var (
	ErrMissingDatabaseURL = errors.New("DATABASE_URL is required")
	ErrMissingJWTSecret   = errors.New("JWT_SECRET is required")
)
