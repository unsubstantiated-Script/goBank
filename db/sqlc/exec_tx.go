package db

import (
	"context"
	"fmt"
)

// execTx executes a function within a database transaction not exported to keep it safe, exported in the function below
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	//Returning a query object
	q := New(tx)

	//Running the functional query
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}
	return tx.Commit(ctx)
}
