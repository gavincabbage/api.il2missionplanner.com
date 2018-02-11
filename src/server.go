package main

import (
	"context"
	"github.com/gavincabbage/api.il2missionplanner.com/src/config"
	"github.com/gorilla/mux"
	"net/http"
)

var version string

type Server struct {
	router *mux.Router
	config *config.Config
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Set some headers to allow requests from anywhere. If we had any sort of user authentication going on
	// this might open us up to CSRF etc. but we don't, and in general I intend for this API to be public.
	// TODO When we pull the mission streaming functionality into this API using websockets, we need to think
	// TODO more carefully about origins there, because then there is some (unimportant) user auth going on.
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	if r.Method == "OPTIONS" {
		return
	}

	if version != "" {
		w.Header().Set("api.il2missionplanner.com.version", version)
	}
	w.Header().Set("api.il2missionplanner.com.env", s.config.Env)

	s.router.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "config", s.config)))
}
