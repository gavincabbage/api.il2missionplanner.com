package sharing

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait           = 10 * time.Second
	pongWait            = 60 * time.Second
	pingPeriod          = (pongWait * 9) / 10
	maxMessageSize      = 2056
	websocketBufferSize = 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  websocketBufferSize,
	WriteBufferSize: websocketBufferSize,
}

// Client is a middleman between the websocket connection and the Room.
type Client struct {
	Room *Room
	Conn *websocket.Conn
	Send chan []byte
}

// ReadMessages pumps messages from the websocket connection to the Room.
func (c *Client) ReadMessages() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()
	c.configureWebsocketConn()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("websocket error: %v", err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.Room.Broadcast <- message
	}
}

func (c *Client) configureWebsocketConn() {
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

// WriteMessages pumps messages from the Room to the websocket connection.
func (c *Client) WriteMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The Room closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
