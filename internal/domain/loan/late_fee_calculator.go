package loan

import (
	"time"
)

const (
	LateFeePerDay = 1.0 // $1 per day late fee
)

// LateFeeCalculator - Domain Service for calculating late fees
type LateFeeCalculator struct{}

func NewLateFeeCalculator() *LateFeeCalculator {
	return &LateFeeCalculator{}
}

// CalculateLateFee calculates the late fee based on return date and due date
// Returns 0.0 if returned on time or before due date
func (c *LateFeeCalculator) CalculateLateFee(dueDate time.Time, returnDate time.Time) float64 {
	// If returned on time or early, no late fee
	if returnDate.Before(dueDate) || returnDate.Equal(dueDate) {
		return 0.0
	}

	// Calculate days late
	daysLate := int(returnDate.Sub(dueDate).Hours() / 24)
	if daysLate < 0 {
		daysLate = 0
	}

	// Calculate late fee
	return float64(daysLate) * LateFeePerDay
}

// IsOverdue checks if a loan is overdue given the current time
func (c *LateFeeCalculator) IsOverdue(dueDate time.Time, now time.Time) bool {
	return now.After(dueDate)
}

// GetDaysLate returns the number of days late (0 if not late)
func (c *LateFeeCalculator) GetDaysLate(dueDate time.Time, returnDate time.Time) int {
	if returnDate.Before(dueDate) || returnDate.Equal(dueDate) {
		return 0
	}

	daysLate := int(returnDate.Sub(dueDate).Hours() / 24)
	if daysLate < 0 {
		return 0
	}

	return daysLate
}
