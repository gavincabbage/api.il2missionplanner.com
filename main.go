package main

import (
	"flag"
	"github.com/gavincabbage/api.il2missionplanner.com/config"
	"github.com/gavincabbage/api.il2missionplanner.com/handlers"
	"github.com/gavincabbage/api.il2missionplanner.com/server"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	configFilePath := flag.String("conf", "conf/conf.json", "path to json configuration file")
	flag.Parse()

	config := config.FromFile(configFilePath)
	log.Println("Host:", config.Host, "Port:", config.Port)

	router := mux.NewRouter()
	router.HandleFunc("/health", handlers.HealthHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/config", handlers.ConfigHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/servers/{server}", handlers.ServerStateHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/servers", handlers.ServersHandler).Methods("GET", "OPTIONS")

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	http.Handle("/", &server.Server{router, config})
	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, nil))
}
