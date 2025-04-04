// TransferTX выполняет денежный перевод с одного счета на другой
// Он создает transfer, add account entries и обновляeт баланс учетных записей в рамках транзакции базы данных

package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// Структура предоставляющая все функции для выполнения отдельных запросов к базе данных и их комбинации в рамках транзакции
// В структуре Queries каждый запрос выполняет только одну операцию с одной конкретной таблицей
// Структура Queries не поддерживает транзакции
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new Store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
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
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	// Создание объекта для выполнения заранее подготовленных SQL-запросов (обёртку для SQL-запросов)
	// Далее через q можно выполнять SQL-запросы внутри транзакции
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
