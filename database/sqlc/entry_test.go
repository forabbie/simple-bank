package db

import (
	"context"
	"testing"
	"time"

	"github.com/forabbie/simplebank/util"
	"github.com/stretchr/testify/require"
)


func TestCreatEntry(t *testing.T) {
	account := createRandomAccount(t)
	arg := CreateEntryParams {
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	arg := CreateEntryParams {
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}
	entry1, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry1)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	arg := CreateEntryParams {
		AccountID: account.ID,
		Amount: util.RandomMoney(),
	}
	for i := 0; i < 10; i++ {
		_, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
	}

	listArg := ListEntriesParams {
		AccountID: account.ID,
		Limit: 5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntries(context.Background(), listArg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}