package websocketserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second
	// Time allowd to read the next pong message from the peer
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10
	// maximum message size allowed from peer
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// needed to allow connections from any origin for :3000 -> :5555
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketServer struct {
}

type JsonData struct {
	Username  string `json:"username"`
	Room      string `json:"room"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type Client struct {
	username string
	room     string
	hub      *ConnHub
	conn     *websocket.Conn
	send     chan JsonData
}

func (c *Client) ReadPump() {
	// schedule client to be disconnected
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// init Client connection
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// handle connection read
	for {
		fmt.Println("reading from client")
		// read JSON data from connection
		message := JsonData{}
		if err := c.conn.ReadJSON(&message); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		fmt.Printf("Got response %#v\n", message)

		// queue messge for writing
		c.hub.send <- message
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		fmt.Println("Sent")
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// channel has been closed by the hub
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			messageJson, _ := json.Marshal(message)
			w.Write(messageJson)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				messageJson, _ := json.Marshal(<-c.send)
				w.Write(messageJson)
			}

			if err := w.Close(); err != nil {
				return
			}

		// send ping over websocket
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (ws WebSocketServer) WsHandler(hub *ConnHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	queries := r.URL.Query()

	// init new client, register to hub
	username := queries.Get("username")
	room := queries.Get("room")

	client := &Client{
		username: username,
		room:     room,
		conn:     conn,
		hub:      hub,
		send:     make(chan JsonData, 256),
	}
	client.hub.register <- client

	go client.WritePump()
	go client.ReadPump()

}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{}
}
