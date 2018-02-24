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
	roomName := mux.Vars(r)["roomName"]
	if hub, exists := (*hubs)[roomName]; exists {
		serveWs(hub, w, r)
	} else {
		room := sharing.NewRoom()
		go room.Start()
		(*hubs)[roomName] = room
		serveWs(room, w, r)
	}
}

func serveWs(room *sharing.Room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &sharing.Client{Room: room, Conn: conn, Send: make(chan []byte, 256)}
	client.Room.Register <- client
	go client.WriteMessages()
	client.ReadMessages()
}

func fromRequestContext(r *http.Request) *map[string]*sharing.Room {
	return r.Context().Value("hubs").(*map[string]*sharing.Room)
}

func SharingHtmlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "./sharing/home.html")
}
