package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Iowel/course-simple-bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Создание подключения к БД // conn - Объект подключения
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = NewStore(connPool)

	// Запуск тестов
	code := m.Run()

	// Завершение работы с кодом возврата
	os.Exit(code)
}
