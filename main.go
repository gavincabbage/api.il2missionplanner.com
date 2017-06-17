package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var servers = map[string]string{
    "randomexpert": "http://72ag-ded.ru/static/il2missionplanner.json",
    "randomnormal": "http://72ag-ded.xyz/static/il2missionplanner.json",
    "virtualpilotsfi": "http://ts3.virtualpilots.fi/output.json",
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func serveServerState(w http.ResponseWriter, r *http.Request) {
    server := mux.Vars(r)["server"]
    url := servers[server]

	response, err := http.Get(url);
    if err != nil {
		log.Fatal("ERROR getting random expert data: ", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body);
    if err != nil {
		log.Fatal("ERROR reading random expert response body", err)
	}
	w.Write(body)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", serveHome).Methods("GET")
	router.HandleFunc("/servers/{server}", serveServerState).Methods("GET")
	log.Fatal(http.ListenAndServe(":9090", router))
}
