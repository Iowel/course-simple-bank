package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Iowel/course-simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	HashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: HashedPassword,
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

func TestUpdateUserOnlyFullName(t *testing.T) {
	// Создадим старого пользователя
	oldUser := createRandomUser(t)

	// Генерируем новое полное имя
	newFullName := util.RandomOwner()
	// Обновим пользователя с именем existing_username согласно переданным данным.
	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})

	// Проверяем что ошибка не равна нулю
	require.NoError(t, err)

	// Проверяем что полное имя обновленного пользователя не совпадает со старым
	require.NotEqual(t, oldUser.FullName, updateUser.FullName)

	// Проверяем что полное имя обновленного пользователя совпадает с новым
	require.Equal(t, newFullName, updateUser.FullName)

	// Проверяем что почта и пароль не изменились
	require.Equal(t, oldUser.Email, updateUser.Email)
	require.Equal(t, oldUser.HashedPassword, updateUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()

	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	// Проверяем что ошибка не равна нулю
	require.NoError(t, err)

	// Проверяем что старая почта не совпадает с новой измененной
	require.NotEqual(t, oldUser.Email, updateUser.Email)

	require.Equal(t, newEmail, updateUser.Email)

	require.Equal(t, oldUser.FullName, updateUser.FullName)
	require.Equal(t, oldUser.HashedPassword, updateUser.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	// Проверяем что ошибка не равна нулю
	require.NoError(t, err)

	require.NotEqual(t, oldUser.HashedPassword, updateUser.HashedPassword)

	require.Equal(t, newHashedPassword, updateUser.HashedPassword)

	require.Equal(t, oldUser.FullName, updateUser.FullName)
	require.Equal(t, oldUser.Email, updateUser.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(6)

	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	// Проверяем что ошибка не равна нулю
	require.NoError(t, err)

	require.NotEqual(t, oldUser.HashedPassword, updateUser.HashedPassword)
	require.Equal(t, newHashedPassword, updateUser.HashedPassword)

	require.NotEqual(t, oldUser.Email, updateUser.Email)
	require.Equal(t, newEmail, updateUser.Email)

	require.NotEqual(t, oldUser.FullName, updateUser.FullName)
	require.Equal(t, newFullName, updateUser.FullName)
}
