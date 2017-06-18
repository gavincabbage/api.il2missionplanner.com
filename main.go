package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"encoding/json"
	"flag"
)

var servers = map[string]string{
	"randomexpert":    "http://72ag-ded.ru/static/il2missionplanner.json",
	"randomnormal":    "http://72ag-ded.xyz/static/il2missionplanner.json",
	"virtualpilotsfi": "http://ts3.virtualpilots.fi/output.json",
}

type Config struct {
	Host    string            `json:"host"`
	Port    string            `json:"port"`
	Servers map[string]string `json:"servers"`
}

type ApiResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func sliceOfMapKeys(m map[string]string) []string {
	slice := make([]string, len(servers))
	i := 0
	for k := range m {
		slice[i] = k
		i++
	}
	return slice
}

func marshalApiResponse(data interface{}, errorMessage string) []byte {
	response := &ApiResponse{data, errorMessage}
	marshalledResponse, err := json.Marshal(response)
	if err != nil {
		marshalledResponse = []byte{}
	}
	return marshalledResponse
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(marshalApiResponse("OK", ""))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, string(marshalApiResponse("", "not found")), 404)
}

func serversHandler(w http.ResponseWriter, r *http.Request) {
	serverList := sliceOfMapKeys(servers)
	w.Write(marshalApiResponse(serverList, ""))
}

func serverStateHandler(w http.ResponseWriter, r *http.Request) {
	server := mux.Vars(r)["server"]
	url := servers[server]
	if url == "" {
		http.Error(w, string(marshalApiResponse("", "server not found")), 404)
		return
	}

	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil || response.StatusCode != 200 {
		http.Error(w, string(marshalApiResponse("", "problem retrieving response from server")), 500)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, string(marshalApiResponse("", "problem reading response from server")), 500)
		return
	}

	var unmarshalledBody map[string]interface{}
	err = json.Unmarshal(body, &unmarshalledBody)
	if err != nil {
		http.Error(w, string(marshalApiResponse("", "problem parsing response from server")), 500)
		return
	}

	w.Write(marshalApiResponse(unmarshalledBody, ""))
}

func main() {

	configFilePath := flag.String("conf", "conf/conf.json", "path to json configuration file")
	flag.Parse()

	rawFileContent, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	config := &Config{}
	err = json.Unmarshal(rawFileContent, config)
	if err != nil {
		log.Fatal(err.Error())
	}
	hostAndPort := config.Host + ":" + config.Port

	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/servers/{server}", serverStateHandler).Methods("GET")
	router.HandleFunc("/servers", serversHandler).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	log.Fatal(http.ListenAndServe(hostAndPort, router))
}
