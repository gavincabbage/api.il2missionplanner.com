package handlers

import (
	"encoding/json"
	"github.com/gavincabbage/api.il2missionplanner.com/config"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(marshalApiResponse("OK", ""))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, string(marshalApiResponse("", "not found")), http.StatusNotFound)
}

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(marshalApiResponse(config.FromRequestContext(r), ""))
}

func marshalApiResponse(data interface{}, errorMessage string) []byte {
	response := make(map[string]interface{})
	response["data"], response["error"] = data, errorMessage
	marshalledResponse, err := json.Marshal(response)
	if err != nil {
		marshalledResponse = []byte{}
	}
	return marshalledResponse
}
