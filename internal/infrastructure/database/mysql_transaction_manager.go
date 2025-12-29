package database

import (
	"database/sql"
)

// MySQLTransactionManager implements TransactionManager for MySQL
type MySQLTransactionManager struct {
	db *sql.DB
}

func NewMySQLTransactionManager(db *sql.DB) *MySQLTransactionManager {
	return &MySQLTransactionManager{db: db}
}

func (tm *MySQLTransactionManager) Transaction(callback func() error) error {
	// Step 1: Begin transaction
	tx, err := tm.db.Begin()
	if err != nil {
		return err
	}

	// Ensure rollback or commit is called
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after rollback
		}
	}()

	// Step 2: Execute callback
	if err := callback(); err != nil {
		// Step 4: Rollback on any error
		tx.Rollback()
		return err
	}

	// Step 3: Commit if successful
	return tx.Commit()
}