package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]any

// WriteJSON serialises data as indented JSON and writes it with the given HTTP status.
func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("WriteJSON: %w", err)
	}
	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

// ReadParamID reads the {id} URL parameter from a chi route.
func ReadParamID(r *http.Request) (string, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return "", errors.New("invalid url param")
	}
	return id, nil
}

// PaginationParams holds validated pagination query parameters.
type PaginationParams struct {
	Page  int
	Limit int
}

// Offset returns the SQL OFFSET value derived from Page and Limit.
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

// ReadPaginationParams reads ?page= and ?limit= from the request URL,
// applying safe defaults and upper bounds.
func ReadPaginationParams(r *http.Request) PaginationParams {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return PaginationParams{Page: page, Limit: limit}
}

// ProductFilter holds optional query-string filters for product list endpoints.
type ProductFilter struct {
	Name  string
	City  string
	State string
}

// ReadProductFilter reads ?q=, ?city= and ?state= from the request URL.
func ReadProductFilter(r *http.Request) ProductFilter {
	return ProductFilter{
		Name:  r.URL.Query().Get("q"),
		City:  r.URL.Query().Get("city"),
		State: r.URL.Query().Get("state"),
	}
}

// RFQFilter holds optional query-string filters for RFQ list endpoints.
type RFQFilter struct {
	ProductName string
	City        string
	State       string
}

// ReadRFQFilter reads ?q=, ?city= and ?state= from the request URL.
func ReadRFQFilter(r *http.Request) RFQFilter {
	return RFQFilter{
		ProductName: r.URL.Query().Get("q"),
		City:        r.URL.Query().Get("city"),
		State:       r.URL.Query().Get("state"),
	}
}
