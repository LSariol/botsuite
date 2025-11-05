package api

import "net/http"

func (s *APIServer) defineRoutes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)

	return mux
}
