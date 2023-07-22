package db

import (
	"context"
	"database/sql"
	"fmt"
)

/**
** Read Phenomena **
Dirty Read: transaction reads data written by another concurrent uncommitted transaction
* this is bad because the transaction might read from the other transaction, and it can do anything for a while before it's closed.
Non-Repeatable Read: transaction reads same row twice and sees different values because it has been modified by the other committed transaction.
Phantom Read: transaction is re-executing a query to find rows that satisfy a condition and sees a different set of rows, due to changes by another committted transaction.
Serialization Anomaly: the result of a group of concurrent committed transactions is impossible to achieve if we try to run them sequentially in any order w/o overlapping.

** 4 Standard Isolation Levels weakest to strongest controls **
1. Read Uncommitted: Can see data written by uncommitted transaction -> allows dirty read.
2. Read Committed: Only see data written by committed transaction -> doesn't allow dirty read.
3. Repeatable Read: Same read query always returns the same result
4. Serializable: Can achieve same result if execute transactions serially in some order instead of concurrently.
*/

// Store provides all functions to execute db queries and transactions
// Embedding Queries here allows store to have access to all the Queries methods and such
type Store struct {
	//Composition here over inheritance.
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction not exported to keep it safe, exported in the function below
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	//Returning a query object
	q := New(tx)

	//Running the functional query
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// This method performs a money transfer from one account to another
// It creates a transfer record, add account entries, and updates accounts' balance within a single DB transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		//Updating this to swip/swap account updates based upon account ID changes to avoid Deadlock
		if arg.FromAccountID < arg.ToAccountID {
			//Go can load three vars at once if three returns are given!
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	//Adding money to account 1
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})

	if err != nil {
		//Go doesn't need to spell out the returns here, it will just do it.
		return
	}

	//Adding money to account 2
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	if err != nil {
		//Go doesn't need to spell out the returns here, it will just do it.
		return
	}

	//Go doesn't need to spell out the returns here, it will just do it.
	return
}
