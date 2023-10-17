package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferWithTx(t *testing.T) {
	store := NewSQLStore(dbConn)

	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	fromAccountBalance := fromAccount.Balance
	toAccountBalance := toAccount.Balance

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	occurrences := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, fromAccount.ID)
		require.Equal(t, transfer.ToAccountID, toAccount.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, fromAccount.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, toAccount.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, fromAccount.ID, result.FromAccount.ID)

		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, toAccount.ID, result.ToAccount.ID)

		diffFromAccount := fromAccountBalance - result.FromAccount.Balance
		diffToAccount := result.ToAccount.Balance - toAccountBalance
		require.Equal(t, diffFromAccount, diffToAccount)
		require.True(t, diffFromAccount > 0)
		require.True(t, diffFromAccount%amount == 0)

		k := int(diffFromAccount / amount)
		require.True(t, k >= 1 && k <= n)

		require.NotContains(t, occurrences, k)
		occurrences[k] = true
	}

	updateFromAccount, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateFromAccount)

	updateToAccount, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateToAccount)

	require.Equal(t, fromAccountBalance-int64(n)*amount, updateFromAccount.Balance)
	require.Equal(t, toAccountBalance+int64(n)*amount, updateToAccount.Balance)

}

func TestTransferWithTxDeadlock(t *testing.T) {
	store := NewSQLStore(dbConn)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount1)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount2)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)

}
