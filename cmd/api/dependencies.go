package main

import (
	"database/sql"

	"library-management/internal/application/getbook"
	"library-management/internal/http/controllers"
	queryImpl "library-management/internal/infrastructure/query"
)

// Dependencies holds all application dependencies
type Dependencies struct {
	BookController *controllers.BookController
	// Add other controllers/services here
}

// SetupDependencies wires up all dependencies
func SetupDependencies(db *sql.DB) *Dependencies {
	// Register QueryService (CQRS Query Side)
	bookQueryService := queryImpl.NewBookQueryService(db)

	// Register use cases
	getBookUseCase := getbook.NewGetBookUseCase(bookQueryService)

	// Register controllers
	bookController := controllers.NewBookController(getBookUseCase)

	return &Dependencies{
		BookController: bookController,
	}
}

// Example usage in main.go:
// db := connectToDatabase()
// deps := SetupDependencies(db)
// setupRoutes(router, deps.BookController)