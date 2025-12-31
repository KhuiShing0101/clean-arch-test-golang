package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"library-management/internal/application/getbook"
)

type BookController struct {
	getBookUseCase *getbook.GetBookUseCase
}

func NewBookController(getBookUseCase *getbook.GetBookUseCase) *BookController {
	return &BookController{
		getBookUseCase: getBookUseCase,
	}
}

// GetBook handles GET /api/books/{bookId}
func (c *BookController) GetBook(w http.ResponseWriter, r *http.Request) {
	// 1. Extract bookId from URL path
	vars := mux.Vars(r)
	bookId := vars["bookId"]

	// 2. Create request DTO
	request, err := getbook.NewGetBookRequest(bookId)
	if err != nil {
		// Request DTO validation failed
		c.sendJSON(w, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// 3. Execute use case
	response := c.getBookUseCase.Execute(request)

	// 4. Handle failure responses
	if !response.Success {
		c.handleError(w, *response.ErrorMessage)
		return
	}

	// 5. Return success response
	c.sendJSON(w, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":           response.Data.ID,
			"isbn":         response.Data.ISBN,
			"title":        response.Data.Title,
			"author":       response.Data.Author,
			"isAvailable":  response.Data.IsAvailable,
			"currentLoan":  response.Data.CurrentLoan,
		},
	}, http.StatusOK)
}

// handleError maps error messages to appropriate HTTP status codes
func (c *BookController) handleError(w http.ResponseWriter, errorMessage string) {
	statusCode := http.StatusInternalServerError

	if strings.Contains(strings.ToLower(errorMessage), "not found") {
		statusCode = http.StatusNotFound
	} else if strings.Contains(errorMessage, "Invalid") {
		statusCode = http.StatusBadRequest
	}

	c.sendJSON(w, map[string]interface{}{
		"success": false,
		"error":   errorMessage,
	}, statusCode)
}

// sendJSON is a helper to send JSON responses
func (c *BookController) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}