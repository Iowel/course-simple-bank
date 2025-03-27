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

// Параметры для транзакции
// Структура содержит все необходимые входные данные для перевода между двумя счетами
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// Структура результата транзакции, содержит инфомрацию о транзакции
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`     // Запись о переводе
	FromAccount Account  `json:"from_account"` // Изменение баланса на счете
	ToAccount   Account  `json:"to_account"`   // Изменение баланса на счете на который были отправлены деньги
	FromEntry   Entry    `json:"from_enrty"`   // Запись о том, что деньги были отправлены со счета
	ToEntry     Entry    `json:"to_entry"`     // Запись о том, что деньги были зачислены на счет
}

// Функция для перевода денег из одного аккаунта в другой
// Она создаст новую запись о переводе, добавит новые записи об аккаунте и обновит баланс аккаунта в рамках одной транзакции с БД
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	// Создание записи о переводе
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Добавление записи об аккаунте для отправителя
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// Добавление записи об аккаунте для получателя
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update accounts balance
		// Получаем запись о счете из БД
		// Добавляем или вычитываем некоторую сумму
		// Обновляем запись в БД

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
