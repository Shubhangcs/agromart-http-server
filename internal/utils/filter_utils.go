package utils

import "net/http"

type ProductFilter struct {
	Name  string
	City  string
	State string
}

func ReadProductFilter(r *http.Request) ProductFilter {
	return ProductFilter{
		Name:  r.URL.Query().Get("q"),
		City:  r.URL.Query().Get("city"),
		State: r.URL.Query().Get("state"),
	}
}

type RFQFilter struct {
	ProductName string
	City        string
	State       string
}

func ReadRFQFilter(r *http.Request) RFQFilter {
	return RFQFilter{
		ProductName: r.URL.Query().Get("q"),
		City:        r.URL.Query().Get("city"),
		State:       r.URL.Query().Get("state"),
	}
}
