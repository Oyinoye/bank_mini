package db

import (
	"context"
	"fmt"
	"math/big"
)

// Money helpers: sqlc maps DECIMAL to string, so balances and amounts are not numeric
// types in Go—we must parse, operate with math/big, and format back to strings for SQL.

// parseDecimalRat parses a database decimal string into *big.Rat for exact arithmetic.
func parseDecimalRat(s string) (*big.Rat, error) {
	r := new(big.Rat)
	if _, ok := r.SetString(s); !ok {
		return nil, fmt.Errorf("invalid decimal amount %q", s)
	}
	return r, nil
}

// ratToDecimalString renders a *big.Rat as a decimal string suitable for query parameters.
func ratToDecimalString(r *big.Rat) string {
	if r == nil {
		return "0"
	}
	f := new(big.Float).SetPrec(256).SetRat(r)
	return f.Text('f', -1)
}

// addMoneyStrings returns the sum of two decimal money strings (a + b).
func addMoneyStrings(a, b string) (string, error) {
	ra, err := parseDecimalRat(a)
	if err != nil {
		return "", err
	}
	rb, err := parseDecimalRat(b)
	if err != nil {
		return "", err
	}
	sum := new(big.Rat).Add(ra, rb)
	return ratToDecimalString(sum), nil
}

// subMoneyStrings returns the difference of two decimal money strings (a - b).
func subMoneyStrings(a, b string) (string, error) {
	ra, err := parseDecimalRat(a)
	if err != nil {
		return "", err
	}
	rb, err := parseDecimalRat(b)
	if err != nil {
		return "", err
	}
	diff := new(big.Rat).Sub(ra, rb)
	return ratToDecimalString(diff), nil
}

// mulMoneyByInt64 multiplies a decimal money string by an integer (e.g. n concurrent transfers).
func mulMoneyByInt64(amount string, n int64) (string, error) {
	ra, err := parseDecimalRat(amount)
	if err != nil {
		return "", err
	}
	p := new(big.Rat).Mul(ra, big.NewRat(n, 1))
	return ratToDecimalString(p), nil
}


// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        string `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction.
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
    var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

        // txName := ctx.Value(txKey)  // cos txKey was added to context in store_test.go

        // fmt.Println(txName, "create transfer")

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

        // fmt.Println(txName, "create entry 1")

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    fmt.Sprintf("-%s", arg.Amount),
		})
		if err != nil {
			return err
		}

        // fmt.Println(txName, "create entry 2")

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

        // Lock and update accounts in a deterministic ID order to avoid deadlocks.
        firstID := arg.FromAccountID
        secondID := arg.ToAccountID
        if firstID > secondID {
            firstID, secondID = secondID, firstID
        }

        firstAccount, err := q.GetAccountForUpdate(ctx, firstID)
        if err != nil {
            return err
        }

        secondAccount, err := q.GetAccountForUpdate(ctx, secondID)
        if err != nil {
            return err
        }

        accountsByID := map[int64]Account{
            firstID:  firstAccount,
            secondID: secondAccount,
        }

        fromCurrent := accountsByID[arg.FromAccountID]
        toCurrent := accountsByID[arg.ToAccountID]

        newFromBal, serr := subMoneyStrings(fromCurrent.Balance, arg.Amount)
        if serr != nil {
            return serr
        }
        newToBal, serr := addMoneyStrings(toCurrent.Balance, arg.Amount)
        if serr != nil {
            return serr
        }

        newBalanceByID := map[int64]string{
            arg.FromAccountID: newFromBal,
            arg.ToAccountID:   newToBal,
        }

        firstUpdated, err := q.UpdateAccount(ctx, UpdateAccountParams{
            ID:      firstID,
            Balance: newBalanceByID[firstID],
        })
        if err != nil {
            return err
        }

        secondUpdated, err := q.UpdateAccount(ctx, UpdateAccountParams{
            ID:      secondID,
            Balance: newBalanceByID[secondID],
        })
        if err != nil {
            return err
        }

        updatedByID := map[int64]Account{
            firstID:  firstUpdated,
            secondID: secondUpdated,
        }
        result.FromAccount = updatedByID[arg.FromAccountID]
        result.ToAccount = updatedByID[arg.ToAccountID]


		// if arg.FromAccountID < arg.ToAccountID {
		// 	result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		// } else {
		// 	result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		// }

		return nil
	})

	return result, err
}
