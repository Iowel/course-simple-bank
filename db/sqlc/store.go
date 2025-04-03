// TransferTX выполняет денежный перевод с одного счета на другой
// Он создает transfer, add account entries и обновляeт баланс учетных записей в рамках транзакции базы данных

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

// Структура предоставляющая все функции для выполнения отдельных запросов к базе данных и их комбинации в рамках транзакции
// В структуре Queries каждый запрос выполняет только одну операцию с одной конкретной таблицей
// Структура Queries не поддерживает транзакции
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// Функция для выполнения общей транзакции базы данных
// В качестве входных данных используется контекст и функцию обратного вызова
// Затем запускается новая транзакция базы данных
// Создается новый объект запроса для этой транзакции
// И вызывается функция обратного вызова с созданными запросами
// Затем фиксируем или отменяем транзакцию в зависимости от ошибки возвращенной этой функцией
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Hачалo транзакции в базе данных
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Создание объекта для выполнения заранее подготовленных SQL-запросов (обёртку для SQL-запросов)
	// Далее через q можно выполнять SQL-запросы внутри транзакции
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
