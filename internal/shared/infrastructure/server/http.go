package infrastructure

import (
	"fmt"
	"log"
	"net/http"
)

type HttpServer struct {
	server    *http.Server
	serverMux *http.ServeMux
}

func NewHttpServer(port string) *HttpServer {
	newServer := &HttpServer{}

	newServer.server = &http.Server{Addr: fmt.Sprintf(":%s", port)}
	newServer.serverMux = &http.ServeMux{}

	newServer.server.Handler = newServer.serverMux

	return newServer
}

func (s *HttpServer) RegisterRoute(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.serverMux.HandleFunc(pattern, handler)
}

func (s *HttpServer) Start() {
	log.Printf("Listening on http://localhost%s", s.server.Addr)

	log.Fatal(s.server.ListenAndServe())
}
