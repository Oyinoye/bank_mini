package db

import (
	"context"
	"testing"
)

// createRandomAccount is a test helper that creates a single account in the test DB.
func createRandomAccount(t *testing.T) Account {
	t.Helper()

	arg := CreateAccountParams{
		Owner:    "tom",
		Balance:  "100", // Balance is a string type (mapped from DECIMAL)
		Currency: "USD",
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	if account.Owner != arg.Owner {
		t.Fatalf("unexpected owner: got %s, want %s", account.Owner, arg.Owner)
	}

	if account.Currency != arg.Currency {
		t.Fatalf("unexpected currency: got %s, want %s", account.Currency, arg.Currency)
	}

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}
