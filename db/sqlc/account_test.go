package db

import (
	"context"
	"testing"

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
		Owner:    "tom",
		Balance:  "100", // Balance is a string type (mapped from DECIMAL)
		Currency: "USD",
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
	createRandomAccount(t)
}
