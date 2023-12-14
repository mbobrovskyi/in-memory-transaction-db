package database_test

import (
	"github.com/mbobrovskyi/in-memory-transaction-db/src/database"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInMemoryDatabase_TestForCommitTransaction(t *testing.T) {
	db := database.NewInMemoryDatabase()
	db.Set("key1", "value1")
	db.StartTransaction()
	db.Set("key1", "value2")
	db.Commit()
	value := db.Get("key1")
	require.EqualValues(t, "value2", value)
}

func TestInMemoryDatabase_TestForRollback(t *testing.T) {
	db := database.NewInMemoryDatabase()
	db.Set("key1", "value1")
	db.StartTransaction()
	value := db.Get("key1")
	require.Equal(t, "value1", value)
	db.Set("key1", "value2")
	value = db.Get("key1")
	require.Equal(t, "value2", value)
	db.Rollback()
	value = db.Get("key1")
	require.Equal(t, "value1", value)
}

func TestInMemoryDatabase_NestedTransactions(t *testing.T) {
	db := database.NewInMemoryDatabase()
	db.Set("key1", "value1")
	db.StartTransaction()
	db.Set("key1", "value2")
	value := db.Get("key1")
	require.Equal(t, "value2", value)
	db.StartTransaction()
	value = db.Get("key1")
	require.Equal(t, "value2", value)
	db.Delete("key1")
	db.Commit()
	value = db.Get("key1")
	require.Nil(t, value)
}

func TestInMemoryDatabase_NestedTransactionsWithRollBack(t *testing.T) {
	db := database.NewInMemoryDatabase()
	db.Set("key1", "value1")
	db.StartTransaction()
	db.Set("key1", "value2")
	value := db.Get("key1")
	require.Equal(t, "value2", value)
	db.StartTransaction()
	value = db.Get("key1")
	require.Equal(t, "value2", value)
	db.Delete("key1")
	db.Rollback()
	value = db.Get("key1")
	require.Equal(t, value, "value2")
	db.Commit()
	value = db.Get("key1")
	require.Equal(t, value, "value2")
}
