package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// Установка Gin в тестовый режим, чтобы тесты выполнялись без логирования, для более чистого вывода во время тестирования
func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
