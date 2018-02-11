package server

import (
	"context"
	"github.com/gavincabbage/api.il2missionplanner.com/config"
	"github.com/gorilla/mux"
	"net/http"
)

// Version is populated during the build process using -ldflags. See the build script for more details:
// https://github.com/gavincabbage/api.il2missionplanner.com/bin/build.bash
var Version string

type Server struct {
	Router *mux.Router
	Config *config.Config
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Set some headers to allow requests from anywhere. If we had any sort of user authentication going on
	// this might open us up to CSRF etc. but we don't, and in general I intend for this API to be public.
	// TODO When we pull the mission streaming functionality into this API using websockets, we may need to think
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

	if Version != "" {
		w.Header().Set("api.il2missionplanner.com.version", Version)
	}
	w.Header().Set("api.il2missionplanner.com.env", s.Config.Env)

	w.Header().Set("Content-Type", "application/json")

	s.Router.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "config", s.Config)))
}
