package database

// TransactionManager executes operations within a database transaction
type TransactionManager interface {
	// Transaction executes callback within a database transaction
	// Returns error if callback fails (rolls back transaction)
	Transaction(callback func() error) error
}