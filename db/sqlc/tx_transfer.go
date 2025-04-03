package db

import "context"

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
