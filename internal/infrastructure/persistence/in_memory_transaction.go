package persistence

import "library-management/internal/domain/shared"

// InMemoryTransactionManager is a simple pass-through implementation
// For Lesson 4 learning - real transaction logic comes later
type InMemoryTransactionManager struct{}

// NewInMemoryTransactionManager creates a new in-memory transaction manager
func NewInMemoryTransactionManager() shared.TransactionManager {
	return &InMemoryTransactionManager{}
}

// RunInTransaction executes the function directly (no actual transaction for learning)
func (tm *InMemoryTransactionManager) RunInTransaction(fn func() error) error {
	// For Lesson 4: Just execute the function
	// Real implementation would: BEGIN -> Execute -> COMMIT/ROLLBACK
	return fn()
}
