package utils

import (
	"net/http"
	"strconv"
)

type PaginationParams struct {
	Page  int
	Limit int
}

func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

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
