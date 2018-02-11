package handlers

import (
	"encoding/json"
	"github.com/gavincabbage/api.il2missionplanner.com/src/config"
	"github.com/gavincabbage/api.il2missionplanner.com/src/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

// ServersHandler returns the list of servers for which mission planner integration is available.
func ServersHandler(w http.ResponseWriter, r *http.Request) {
	config := config.FromRequestContext(r)
	serverList := sliceOfMapKeys(config.Servers)
	w.Write(marshalApiResponse(serverList, ""))
}

// ServerStateHandler returns the state of a given server.
func ServerStateHandler(w http.ResponseWriter, r *http.Request) {
	config := config.FromRequestContext(r)
	server := mux.Vars(r)["server"]
	url := config.Servers[server]
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

	var unmarshalledBody models.Plan
	err = json.Unmarshal(body, &unmarshalledBody)
	if err != nil {
		http.Error(w, string(marshalApiResponse("", "problem parsing response from server")), 500)
		return
	}

	w.Write(marshalApiResponse(unmarshalledBody, ""))
}

func sliceOfMapKeys(m map[string]string) []string {
	slice := make([]string, len(m))
	i := 0
	for k := range m {
		slice[i] = k
		i++
	}
	return slice
}
