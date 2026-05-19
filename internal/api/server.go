package broker

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

func (s *APIServer) Initilize() *http.ServeMux {
	mux := s.defineRoutes()
	fmt.Println(mux)
	return mux
}

func (s *APIServer) Run(address string, mux *http.ServeMux) {

	log.Fatal(http.ListenAndServe(address, mux))
}
