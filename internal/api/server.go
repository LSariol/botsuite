package api

import (
	"fmt"
	"log"
	"net/http"
)

type APIServer struct {
	Config *Config
}

func NewAPIServer() *APIServer {
	return &APIServer{}
}

func (s *APIServer) Initilize() error {
	mux := s.defineRoutes()
	fmt.Println(mux)
	return nil
}

func (s *APIServer) Run(address string, mux *http.ServeMux) {

	log.Fatal(http.ListenAndServe(address, mux))
}
