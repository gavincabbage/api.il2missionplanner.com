package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	config Config
	router *chi.Mux
}

func NewServer(c Config, r *chi.Mux) *Server {
	return &Server{
		config: c,
		router: r,
	}
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// TODO change this to only allow whitelisted origins to protect against CSRF.
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	if req.Method == "OPTIONS" {
		return
	}
	s.router.ServeHTTP(rw, req)
}
