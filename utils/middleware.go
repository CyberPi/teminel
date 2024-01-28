package utils

import "net/http"

func Use(handler http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, toChain := range middleware {
		handler = toChain(handler)
	}
	return handler
}
