package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MatheusAbdias/go_simple_bank/util"
)

func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := CreateRandomAccount(t)

	fetchedAcoount, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, fetchedAcoount)

	require.Equal(t, account.ID, fetchedAcoount.ID)
	require.Equal(t, account.Owner, fetchedAcoount.Owner)
	require.Equal(t, account.Balance, fetchedAcoount.Balance)
	require.Equal(t, account.Currency, fetchedAcoount.Currency)
	require.WithinDuration(t, account.CreatedAt, fetchedAcoount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := CreateRandomAccount(t)

	arg := AddAccountBalanceParams{ID: account.ID, Amount: util.RandomMoney()}

	fetchedAcoount, err := testQueries.AddAccountBalance(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, fetchedAcoount)

	require.Equal(t, account.ID, fetchedAcoount.ID)
	require.Equal(t, account.Owner, fetchedAcoount.Owner)
	require.Equal(t, account.Currency, fetchedAcoount.Currency)
	require.Equal(t, account.Balance+arg.Amount, fetchedAcoount.Balance)
	require.WithinDuration(t, account.CreatedAt, fetchedAcoount.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account := CreateRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)

	fetchedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, fetchedAccount)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateRandomAccount(t)
	}

	arg := ListAccountsParams{Limit: 2, Offset: 3}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, int(arg.Limit))

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
