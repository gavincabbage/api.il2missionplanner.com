package handlers

import (
	"github.com/gavincabbage/api.il2missionplanner.com/sharing"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SharingHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Origin")
	hubs := fromRequestContext(r)
	room := mux.Vars(r)["room"]
	if hub, exists := (*hubs)[room]; exists {
		serveWs(hub, w, r)
	} else {
		hub := sharing.NewRoom()
		go hub.Run()
		(*hubs)[room] = hub
		serveWs(hub, w, r)
	}
}

func SharingHtmlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "./sharing/home.html")
}

func serveWs(hub *sharing.Room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &sharing.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client
	go client.WritePump()
	client.ReadPump()
}

func fromRequestContext(r *http.Request) *map[string]*sharing.Room {
	return r.Context().Value("hubs").(*map[string]*sharing.Room)
}
