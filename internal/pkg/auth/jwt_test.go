package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTokenRoundTrip(t *testing.T) {
	SetJWTSecret("test-secret")
	want := uuid.New()

	token, err := GenerateToken(want.String())
	require.NoError(t, err)

	got, err := ParseUserIDFromToken(token)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestParseUserIDRejectsUnexpectedAlgorithm(t *testing.T) {
	SetJWTSecret("test-secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.MapClaims{
		"user_id": uuid.NewString(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte("test-secret"))
	require.NoError(t, err)

	_, err = ParseUserIDFromToken(signed)
	require.Error(t, err)
}
