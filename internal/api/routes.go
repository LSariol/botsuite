package broker

import "net/http"

func (s *APIServer) defineRoutes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/letterboxd", letterboxd)

	return mux
}
