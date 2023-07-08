package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"goBank/util"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		//Owner:    "tom",
		//Balance:  100,
		//Currency: "USD",

		//Randomized instead
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
