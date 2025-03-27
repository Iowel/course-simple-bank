package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Структура предоставляющая все функции для выполнения отдельных запросов к базе данных и их комбинации в рамках транзакции

// В структуре Queries каждый запрос выполняет только одну операцию с одной конкретной таблицей
// Структура Queries не поддерживает транзакции
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
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
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
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

		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{ // Обновляем баланс счета
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}
