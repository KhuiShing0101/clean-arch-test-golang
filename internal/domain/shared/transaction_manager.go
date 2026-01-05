package shared

// TransactionManager defines the interface for managing database transactions
// This follows Dependency Inversion Principle - domain defines the contract
type TransactionManager interface {
	// RunInTransaction executes a function within a database transaction
	// If the function returns an error, the transaction is rolled back
	// If the function succeeds, the transaction is committed
	RunInTransaction(fn func() error) error
}
