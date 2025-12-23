# clean-arch-test-golang
library-management/
├── go.mod
├── go.sum
├── domain/
│   ├── book/
│   │   ├── book_id.go
│   │   ├── isbn.go
│   │   ├── book.go
│   │   └── repository.go
│   └── shared/
│       └── errors/
│           └── domain_error.go
└── domain/
    └── book/
        ├── book_id_test.go
        ├── isbn_test.go
        └── book_test.go

// Go uses built-in testing package
// Test files should be named *_test.go
// Example: BookId_test.go, ISBN_test.go

// Run tests with:
// go test ./...

// Run with coverage:
// go test -cover ./...