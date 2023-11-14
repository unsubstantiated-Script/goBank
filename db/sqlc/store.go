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

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// SQLStore provides all functions to execute SQL db queries and transactions
// Embedding Queries here allows store to have access to all the Queries methods and such
type SQLStore struct {
	//Composition here over inheritance.
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction not exported to keep it safe, exported in the function below
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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
