package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T, fromAccount, toAccount Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        fromAccount.Balance,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	CreateRandomTransfer(t, fromAccount, toAccount)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)
	transfer := CreateRandomTransfer(t, fromAccount, toAccount)

	fetchedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, fetchedTransfer)

	require.Equal(t, fetchedTransfer.Amount, transfer.Amount)
	require.Equal(t, fetchedTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, fetchedTransfer.ToAccountID, transfer.ToAccountID)
	require.WithinDuration(t, fetchedTransfer.CreatedAt, transfer.CreatedAt, time.Second)

}

func TestListTransfer(t *testing.T) {
	fromAccount := CreateRandomAccount(t)
	toAccount := CreateRandomAccount(t)

	for i := 0; i < 5; i++ {
		CreateRandomTransfer(t, fromAccount, toAccount)
	}

	arg := ListTransfersParams{
		Limit:         2,
		Offset:        3,
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, int(arg.Limit))

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(
			t,
			transfer.FromAccountID == arg.FromAccountID || transfer.ToAccountID == arg.ToAccountID,
		)
	}
}
