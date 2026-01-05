# Lesson 6: ListBooks Use Case - Pagination Complete ✅

## What Was Created

### 📁 Application Layer

**1. `internal/application/listbooks/list_books_request.go`** - Request DTO
```go
type ListBooksRequest struct {
    Page  int
    Limit int
}

- NewListBooksRequest(page, limit) - Factory with defaults
- GetOffset() - Calculates SQL offset from page/limit
```

**Business Rules:**
- Default page: 1 (if <= 0)
- Default limit: 10 (if <= 0)
- Offset calculation: `(page - 1) * limit`

**2. `internal/application/listbooks/list_books_response.go`** - Response DTO
```go
type ListBooksResponse struct {
    Books      []BookDTO
    Pagination PaginationMetadata
}

type PaginationMetadata struct {
    Page       int
    Limit      int
    Total      int
    TotalPages int
}
```

**3. `internal/application/listbooks/list_books_usecase.go`** - Use Case
```go
func (uc *ListBooksUseCase) Execute(req *ListBooksRequest) (*ListBooksResponse, error)
```

**What It Does:**
1. ✅ Validates page (must be > 0)
2. ✅ Validates limit (must be > 0)
3. ✅ Calculates offset from page/limit
4. ✅ Queries paginated books via QueryService
5. ✅ Returns books with pagination metadata
6. ✅ Handles edge cases (empty list, out of range pages)

**4. `internal/application/listbooks/list_books_usecase_test.go`** - Tests
```
✅ TestListBooksUseCase_Success_FirstPage
✅ TestListBooksUseCase_Success_SecondPage
✅ TestListBooksUseCase_Success_CustomLimit
✅ TestListBooksUseCase_Error_InvalidPageZero
✅ TestListBooksUseCase_Error_InvalidPageNegative
✅ TestListBooksUseCase_Error_InvalidLimitZero
✅ TestListBooksUseCase_Error_InvalidLimitNegative
✅ TestListBooksUseCase_Success_OutOfRangePage
✅ TestListBooksUseCase_Success_EmptyList
✅ TestListBooksUseCase_Error_QueryServiceFailure
```

### 📁 Query Service Extension

**5. `internal/application/query/book_query_service.go`** - Extended Interface
```go
type BookQueryService interface {
    GetBookById(bookId string) (*BookReadModel, error)

    // NEW: Lesson 6
    ListBooks(limit int, offset int) ([]*BookReadModel, int, error)
}
```

**6. `internal/application/getbook/book_dto.go`** - DTO & Adapter
```go
type BookDTO struct {
    BookId      string
    ISBN        string
    Title       string
    Author      string
    IsAvailable bool
}

// Adapter pattern to bridge BookQueryService and IBookQueryService
type BookQueryServiceAdapter struct {
    queryService query.BookQueryService
}
```

## 🧪 Test Results

```bash
=== RUN   TestListBooksUseCase_Success_FirstPage
--- PASS: TestListBooksUseCase_Success_FirstPage (0.00s)
=== RUN   TestListBooksUseCase_Success_SecondPage
--- PASS: TestListBooksUseCase_Success_SecondPage (0.00s)
=== RUN   TestListBooksUseCase_Success_CustomLimit
--- PASS: TestListBooksUseCase_Success_CustomLimit (0.00s)
=== RUN   TestListBooksUseCase_Error_InvalidPageZero
--- PASS: TestListBooksUseCase_Error_InvalidPageZero (0.00s)
=== RUN   TestListBooksUseCase_Error_InvalidPageNegative
--- PASS: TestListBooksUseCase_Error_InvalidPageNegative (0.00s)
=== RUN   TestListBooksUseCase_Error_InvalidLimitZero
--- PASS: TestListBooksUseCase_Error_InvalidLimitZero (0.00s)
=== RUN   TestListBooksUseCase_Error_InvalidLimitNegative
--- PASS: TestListBooksUseCase_Error_InvalidLimitNegative (0.00s)
=== RUN   TestListBooksUseCase_Success_OutOfRangePage
--- PASS: TestListBooksUseCase_Success_OutOfRangePage (0.00s)
=== RUN   TestListBooksUseCase_Success_EmptyList
--- PASS: TestListBooksUseCase_Success_EmptyList (0.00s)
=== RUN   TestListBooksUseCase_Error_QueryServiceFailure
--- PASS: TestListBooksUseCase_Error_QueryServiceFailure (0.00s)
PASS
ok      library-management/internal/application/listbooks     0.700s
```

**ALL 10 TESTS PASSING! ✅**

## 📝 How to Test

### Run Tests
```bash
# Run all ListBooks tests
go test ./internal/application/listbooks/... -v

# Run specific test
go test ./internal/application/listbooks/... -v -run TestListBooksUseCase_Success_FirstPage

# Run all tests in project
go test ./... -v
```

### Create Git Branch & Commit
```bash
# Create feature branch
git checkout -b feature/lesson-6-listbooks

# Add changes
git add .

# Commit
git commit -m "feat(lesson6): implement ListBooks use case with pagination

- Extend BookQueryService with ListBooks(limit, offset) method
- Implement ListBooksUseCase with pagination support
- Add ListBooksRequest and ListBooksResponse DTOs with metadata
- Create BookQueryServiceAdapter for DTO mapping
- Support page/limit parameters with defaults
- Calculate pagination metadata (total pages, current page)
- Handle edge cases (invalid page/limit, empty results)
- Add 10 test scenarios covering all pagination cases"

# Push to remote
git push origin feature/lesson-6-listbooks
```

## 🎯 Key Learning Points

### 1. **Pagination Pattern**
- **Offset/Limit** - Standard SQL pagination
- Offset = `(page - 1) * limit`
- Efficient for datasets up to ~10,000 records
- For larger datasets, consider cursor-based pagination

### 2. **CQRS Query Side**
- Uses `BookQueryService` from Lesson 4
- Read-only operation (no state changes)
- Returns DTOs, not domain entities
- Optimized for reads (can use denormalized data)

### 3. **Response Metadata Pattern**
- Data + Metadata structure
- Metadata includes: `page`, `limit`, `total`, `totalPages`
- Enables UI to build pagination controls
- Provides full context for client

### 4. **Edge Case Handling**
```go
// Invalid inputs -> Error
page <= 0   → ErrInvalidPage
limit <= 0  → ErrInvalidLimit

// Valid but out of range -> Empty list
page > totalPages → []BookDTO{} (empty)
```

### 5. **Single Source of Truth (from Lesson 4)**
- Book availability derived from loans table
- No `book.status` field needed
- Query service performs JOIN to get availability
- Consistency guaranteed by SSOT pattern

### 6. **N+1 Query Problem Prevention**
```sql
-- BAD: N+1 queries (1 for books + N for each book's loans)
SELECT * FROM books;
for each book:
    SELECT * FROM loans WHERE book_id = ?;

-- GOOD: Single JOIN query
SELECT b.*,
       CASE WHEN l.id IS NOT NULL THEN false ELSE true END as is_available
FROM books b
LEFT JOIN loans l ON b.id = l.book_id AND l.returned_at IS NULL
LIMIT ? OFFSET ?;
```

## 📊 Architecture Diagram

```
┌─────────────────────────────────────────┐
│         ListBooksUseCase                │
│  (Application Layer - Thin Coordinator) │
└─────────────────┬───────────────────────┘
                  │
                  ▼
         ┌────────────────┐
         │ BookQueryService │ ← Lesson 4 (Extended)
         │ (CQRS Query)     │
         └────────┬─────────┘
                  │
         ┌────────▼─────────┐
         │   SQL Query      │
         │  (JOIN books +   │
         │   loans table)   │
         └──────────────────┘

Query Flow:
1. Request(page=2, limit=10)
2. Calculate offset = (2-1) * 10 = 10
3. SQL: SELECT ... LIMIT 10 OFFSET 10
4. Get total count: SELECT COUNT(*)
5. Calculate totalPages = ceil(total / limit)
6. Return Response{ books[], pagination{} }
```

## ✅ API Specification

### Endpoint
```
GET /books?page=1&limit=10
```

### Query Parameters
- `page` - Page number (default: 1, must be > 0)
- `limit` - Items per page (default: 10, must be > 0)

### Response
```json
{
  "books": [
    {
      "bookId": "b-12345",
      "isbn": "978-0-13-468599-1",
      "title": "Clean Architecture",
      "author": "Robert C. Martin",
      "isAvailable": true
    },
    {
      "bookId": "b-67890",
      "isbn": "978-0-13-235088-4",
      "title": "Clean Code",
      "author": "Robert C. Martin",
      "isAvailable": false
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "totalPages": 3
  }
}
```

### Error Responses
```json
// Invalid page
{
  "error": "invalid page number"
}

// Invalid limit
{
  "error": "invalid limit"
}
```

## 🔄 Pagination Examples

### Example 1: First Page
```
Request:  page=1, limit=10
Offset:   (1-1) * 10 = 0
SQL:      LIMIT 10 OFFSET 0
Returns:  Books 1-10
```

### Example 2: Second Page
```
Request:  page=2, limit=10
Offset:   (2-1) * 10 = 10
SQL:      LIMIT 10 OFFSET 10
Returns:  Books 11-20
```

### Example 3: Custom Limit
```
Request:  page=1, limit=5
Offset:   (1-1) * 5 = 0
SQL:      LIMIT 5 OFFSET 0
Returns:  Books 1-5
```

### Example 4: Last Page (Partial)
```
Total:    25 books
Request:  page=3, limit=10
Offset:   (3-1) * 10 = 20
SQL:      LIMIT 10 OFFSET 20
Returns:  Books 21-25 (5 books)
```

### Example 5: Out of Range
```
Total:    25 books
Request:  page=10, limit=10
Offset:   (10-1) * 10 = 90
SQL:      LIMIT 10 OFFSET 90
Returns:  [] (empty array)
```

## 🔗 Next Steps

1. **Create Pull Request** on GitHub
2. **Submit PR URL** in the tutorial app (Lesson 6)
3. **AI Review** will check:
   - ✅ ListBooks use case implementation
   - ✅ Pagination logic (offset/limit)
   - ✅ BookQueryService extension
   - ✅ Single Source of Truth (no book.status!)
   - ✅ Response metadata structure
   - ✅ Edge case handling
   - ✅ Test coverage (10 scenarios)
   - ✅ N+1 query prevention

## 📚 Comparison: GetBook vs ListBooks

| Aspect | GetBook (Lesson 4) | ListBooks (Lesson 6) |
|--------|-------------------|---------------------|
| **Pattern** | CQRS Query | CQRS Query |
| **Input** | Single ID | Page + Limit |
| **Output** | Single DTO | Array + Metadata |
| **Query** | WHERE id = ? | LIMIT ? OFFSET ? |
| **Complexity** | Simple | Moderate |
| **Use Case** | View book details | Browse catalog |
| **Response Time** | Fast (~1ms) | Medium (~10ms) |

Both use the same `BookQueryService` interface! 🎯

## ✅ Ready for Submission!

Your Lesson 6 implementation is complete and all tests pass. Follow the steps above to create a PR and submit for AI review.

**Happy coding! 🚀**
