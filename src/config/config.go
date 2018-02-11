package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	Env     string            `json:"env"`
	Host    string            `json:"host"`
	Port    string            `json:"port"`
	Servers map[string]string `json:"servers"`
}

func FromFile(configFilePath *string) *Config {
	rawFileContent, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	config := &Config{}
	err = json.Unmarshal(rawFileContent, config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return config
}

func FromRequestContext(r *http.Request) *Config {
	return r.Context().Value("config").(*Config)
}
