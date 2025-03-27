package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	HashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, HashedPassword)

	err = CheckPassword(password, HashedPassword)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, HashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
