package database

import (
	"maps"
)

type InMemoryDatabase interface {
	// Get the value associated with the given key.
	Get(key string) any

	// Set a key-value pair in the database.
	Set(key string, value any)

	// Delete the key-value pair associated with the given key.
	Delete(key string)

	// StartTransaction start a new transaction.
	//All operations within this transaction are isolated from others.
	StartTransaction()

	// Commit all changes made within the current transaction to the database.
	Commit()

	// Rollback all changes made within the current transaction and discard them.
	Rollback()
}

type store map[string]any

// Since we do not use multithreading, it is most efficient to use a stack, in which we will store history.
type transactions []store

func (t *transactions) IsEmpty() bool {
	return len(*t) == 0
}

func (t *transactions) Push(transactions store) {
	*t = append(*t, transactions)
}

func (t *transactions) Peak() store {
	return (*t)[len(*t)-1]
}

func (t *transactions) Pop() store {
	s := t.Peak()
	*t = (*t)[:len(*t)-1]
	return s
}

type inMemoryDatabase struct {
	store        store
	transactions transactions
}

func (d *inMemoryDatabase) Get(key string) any {
	return d.store[key]
}

func (d *inMemoryDatabase) Set(key string, value any) {
	d.store[key] = value
}

func (d *inMemoryDatabase) Delete(key string) {
	delete(d.store, key)
}

func (d *inMemoryDatabase) StartTransaction() {
	// Save current state on transactions
	d.transactions = append(d.transactions, maps.Clone(d.store))
}

func (d *inMemoryDatabase) Commit() {
	// Skip if we don't have transactions
	if d.transactions.IsEmpty() {
		return
	}

	d.transactions.Pop()
}

func (d *inMemoryDatabase) Rollback() {
	// Skip if we don't have transactions
	if d.transactions.IsEmpty() {
		return
	}

	// Revert previous state
	d.store = d.transactions[len(d.transactions)-1]

	d.transactions.Pop()
}

func NewInMemoryDatabase() InMemoryDatabase {
	return &inMemoryDatabase{
		store: make(store),
	}
}
