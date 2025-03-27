package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Iowel/course-simple-bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Создание подключения к БД // conn - Объект подключения
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	// Запуск тестов
	code := m.Run()

	// Завершение работы с кодом возврата
	os.Exit(code)
}
