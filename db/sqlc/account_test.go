package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Oyinoye/bank_mini/util"
	"github.com/stretchr/testify/require"
)

// // createRandomAccount is a test helper that creates a single account in the test DB.
// func createRandomAccount(t *testing.T) Account {
// 	t.Helper()

// 	arg := CreateAccountParams{
// 		Owner:    "tom",
// 		Balance:  "100", // Balance is a string type (mapped from DECIMAL)
// 		Currency: "USD",
// 	}

// 	account, err := testQueries.CreateAccount(context.Background(), arg)
// 	if err != nil {
// 		t.Fatalf("failed to create account: %v", err)
// 	}

// 	if account.Owner != arg.Owner {
// 		t.Fatalf("unexpected owner: got %s, want %s", account.Owner, arg.Owner)
// 	}

// 	if account.Currency != arg.Currency {
// 		t.Fatalf("unexpected currency: got %s, want %s", account.Currency, arg.Currency)
// 	}

// 	return account
// }

// Using go's testfy package
// createRandomAccount is a test helper that creates a single account in the test DB.
func createRandomAccount(t *testing.T) Account {

	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(), // Balance is a string type (mapped from DECIMAL)
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

func createAccountWithOwner(t *testing.T, owner string) Account {
	t.Helper()

	arg := CreateAccountParams{
		Owner:    owner,
		Balance:  util.RandomMoney(), // Balance is a string type (mapped from DECIMAL)
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}


func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}


func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}


func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	owner := "test-owner"
	for i := 0; i < 10; i++ {
		createAccountWithOwner(t, owner)
	}

	arg := ListAccountsParams{
		Owner:  owner,
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	// require.NotEmpty(t, accounts)

	for _, account := range accounts {
        require.NotEmpty(t, account)
		require.Equal(t, owner, account.Owner)
	}
}
