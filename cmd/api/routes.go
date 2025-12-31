package main

import (
	"github.com/gorilla/mux"
	"library-management/internal/http/controllers"
)

func setupRoutes(r *mux.Router, controller *controllers.BookController) {
	// Route with path parameter
	r.HandleFunc("/books/{bookId}", controller.GetBook).Methods("GET")

	// API prefix
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/books/{bookId}", controller.GetBook).Methods("GET")
}