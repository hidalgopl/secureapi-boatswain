package statusendpoints

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) AttachRoutes(router *mux.Router) {
	//TODO replace with proper status endpoints
	router.HandleFunc("/ping", s.ping()).
		Methods("GET")
}

func (s *Server) ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "pong")
	}
}
