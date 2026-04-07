package db

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	fmt.Println("----Start TestTransfer test----")

    store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

    // number of concurrent transfer transactions
	n := 2
	amount := "10"

    // Use channels to enable data sharing between different go routines.
    // make keyward creates the channel
	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {

	// 	go func() {
	// 		result, err := store.TransferTx(context.Background(), TransferTxParams{
	// 			FromAccountID: account1.ID,
	// 			ToAccountID:   account2.ID,
	// 			Amount:        amount,
	// 		})

    //         // Send errors and results to respective channel.   channel <- data
	// 		errs <- err
	// 		results <- result
	// 	}()

        // Commented - This code below is used for debugging deadlock scenario.
        txName := fmt.Sprintf("tx %d", i+1)
		go func() {
            ctx := context.WithValue(context.Background(), txKey, txName) // now context will hold transaction key and name
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

            // Send errors and results to respective channel.   channel <- data
			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		// From-entry amounts are stored as negative strings (e.g. "-10"), not Go negation on string.
		require.Equal(t, fmt.Sprintf("-%s", amount), fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balances: use package money helpers so DECIMAL strings compare correctly.
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1Str, err := subMoneyStrings(account1.Balance, fromAccount.Balance)
		require.NoError(t, err)
		diff2Str, err := subMoneyStrings(toAccount.Balance, account2.Balance)
		require.NoError(t, err)
		require.Equal(t, diff1Str, diff2Str)

		// Each goroutine debits one transfer amount; diff1 must be k×amount for some k in [1,n].
		diff1Rat, err := parseDecimalRat(diff1Str)
		require.NoError(t, err)
		require.True(t, diff1Rat.Sign() > 0)

		amountRat, err := parseDecimalRat(amount)
		require.NoError(t, err)
		quotient := new(big.Rat).Quo(new(big.Rat).Set(diff1Rat), amountRat)
		require.True(t, quotient.IsInt())

		k := int(quotient.Num().Int64())
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	// Expect total moved = n × amount; compare balances with Rat.Cmp so DB formatting (e.g. scale) does not false-fail.
	totalMoved, err := mulMoneyByInt64(amount, int64(n))
	require.NoError(t, err)
	expBal1, err := subMoneyStrings(account1.Balance, totalMoved)
	require.NoError(t, err)
	expBal2, err := addMoneyStrings(account2.Balance, totalMoved)
	require.NoError(t, err)

	ub1, err := parseDecimalRat(updatedAccount1.Balance)
	require.NoError(t, err)
	eb1, err := parseDecimalRat(expBal1)
	require.NoError(t, err)
	require.Equal(t, 0, eb1.Cmp(ub1))

	ub2, err := parseDecimalRat(updatedAccount2.Balance)
	require.NoError(t, err)
	eb2, err := parseDecimalRat(expBal2)
	require.NoError(t, err)
	require.Equal(t, 0, eb2.Cmp(ub2))
}

// func TestTransferTxDeadlock(t *testing.T) {
// 	account1 := createRandomAccount(t)
// 	account2 := createRandomAccount(t)
// 	fmt.Println(">> before:", account1.Balance, account2.Balance)

// 	n := 10
// 	amount := int64(10)
// 	errs := make(chan error)

// 	for i := 0; i < n; i++ {
// 		fromAccountID := account1.ID
// 		toAccountID := account2.ID

// 		if i%2 == 1 {
// 			fromAccountID = account2.ID
// 			toAccountID = account1.ID
// 		}

// 		go func() {
// 			_, err := store.TransferTx(context.Background(), TransferTxParams{
// 				FromAccountID: fromAccountID,
// 				ToAccountID:   toAccountID,
// 				Amount:        amount,
// 			})

// 			errs <- err
// 		}()
// 	}

// 	for i := 0; i < n; i++ {
// 		err := <-errs
// 		require.NoError(t, err)
// 	}

// 	// check the final updated balance
// 	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
// 	require.NoError(t, err)

// 	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
// 	require.NoError(t, err)

// 	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
// 	require.Equal(t, account1.Balance, updatedAccount1.Balance)
// 	require.Equal(t, account2.Balance, updatedAccount2.Balance)
// }
