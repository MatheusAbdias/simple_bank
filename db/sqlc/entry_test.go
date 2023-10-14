package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MatheusAbdias/go_simple_bank/util"
)

func CreateRandomEntry(t *testing.T, toAccount Account, amount int64) Entry {
	arg := CreateEntryParams{AccountID: toAccount.ID, Amount: amount}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)
	return entry
}

func TestCreateEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	CreateRandomEntry(t, account, util.RandomMoney())
}

func TestGetEntry(t *testing.T) {
	account := CreateRandomAccount(t)
	entry := CreateRandomEntry(t, account, util.RandomMoney())

	fetchedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, fetchedEntry)

	require.Equal(t, fetchedEntry.AccountID, entry.AccountID)
	require.Equal(t, fetchedEntry.Amount, entry.Amount)
	require.WithinDuration(t, fetchedEntry.CreatedAt, entry.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := CreateRandomAccount(t)
	for i := 0; i < 5; i++ {
		CreateRandomEntry(t, account, util.RandomMoney())
	}

	arg := ListEntriesParams{AccountID: account.ID, Limit: 2, Offset: 3}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, int(arg.Limit))

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, account.ID)
	}
}
