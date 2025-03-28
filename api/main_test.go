package api

import (
	"os"
	"testing"
	"time"

	db "github.com/Iowel/course-simple-bank/db/sqlc"
	"github.com/Iowel/course-simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

// Установка Gin в тестовый режим, чтобы тесты выполнялись без логирования, для более чистого вывода во время тестирования
func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
