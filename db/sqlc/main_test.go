package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:1234@localhost:5434/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	// Создание подключения к БД // conn - Объект подключения
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(conn)

	// Запуск тестов
	code := m.Run()

	// Завершение работы с кодом возврата
	os.Exit(code)
}
