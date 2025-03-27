package db

import (
	"context"
	"testing"
	"time"

	"github.com/Iowel/course-simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	// Сравнение полей возвращённого аккаунта с переданными параметрами
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())

	// Проверка на отсутствие нулевых значений
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	// create account
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	// Проверка равенства полей у двух объектов пользователя
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	// Проверка что два времени (например, даты и времени) находятся в пределах указанного диапазона (в данном случае, 1 секунда)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second) // Если разница между PasswordChangedAt пользователей больше 1 секунды, тест завершится ошибкой
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
