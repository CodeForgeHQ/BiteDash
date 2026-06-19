package service

import (
	"os"
	"testing"

	"bitedash/internal/pkg/auth"
)

func TestMain(m *testing.M) {
	auth.SetJWTSecret("test-jwt-secret")
	os.Exit(m.Run())
}
