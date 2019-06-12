package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	api "github.com/gavincabbage/api.il2missionplanner.com"
	"github.com/gavincabbage/api.il2missionplanner.com/log"

	_ "net/http/pprof"
)

func main() {
	logger := logrus.New()

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(kill)
		close(kill)
	}()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-kill
		cancel()
	}()

	var config api.Config
	envconfig.MustProcess("api", &config)

	go func() {
		if err := http.ListenAndServe(":"+config.PprofPort, nil); err != nil {
			logger.Fatal(err)
		}
	}()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(log.Middleware(logger))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	server := api.NewServer(config, router)
	http.Handle("/", server)
	logger.Fatal(http.ListenAndServe(":"+config.Port, router))
}
