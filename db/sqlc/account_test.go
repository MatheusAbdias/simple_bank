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
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
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

	fetchedAccount, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, fetchedAccount)

	require.Equal(t, account.ID, fetchedAccount.ID)
	require.Equal(t, account.Owner, fetchedAccount.Owner)
	require.Equal(t, account.Balance, fetchedAccount.Balance)
	require.Equal(t, account.Currency, fetchedAccount.Currency)
	require.WithinDuration(t, account.CreatedAt, fetchedAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := CreateRandomAccount(t)

	arg := AddAccountBalanceParams{ID: account.ID, Amount: util.RandomMoney()}

	fetchedAccount, err := testQueries.AddAccountBalance(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, fetchedAccount)

	require.Equal(t, account.ID, fetchedAccount.ID)
	require.Equal(t, account.Owner, fetchedAccount.Owner)
	require.Equal(t, account.Currency, fetchedAccount.Currency)
	require.Equal(t, account.Balance+arg.Amount, fetchedAccount.Balance)
	require.WithinDuration(t, account.CreatedAt, fetchedAccount.CreatedAt, time.Second)

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
	var lastAccount Account

	for i := 0; i < 5; i++ {
		lastAccount = CreateRandomAccount(t)
	}

	arg := ListAccountsParams{Owner: lastAccount.Owner, Limit: 5, Offset: 0}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
