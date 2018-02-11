// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	hubs := make(map[string]*Hub)

	router := mux.NewRouter()
	router.HandleFunc("/", serveHome)
	router.HandleFunc("/ws/{room}", func(w http.ResponseWriter, r *http.Request) {
		room := mux.Vars(r)["room"]
		if hub, exists := hubs[room]; exists {
			serveWs(hub, w, r)
		} else {
			hub := newHub()
			go hub.run()
			hubs[room] = hub
			serveWs(hub, w, r)
		}
	})
	http.Handle("/", router)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
