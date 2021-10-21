package websocketserver

import (
	"encoding/json"
	"fmt"
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

type WebSocketServer struct {
	httpport uint
	upgrader websocket.Upgrader
}

type JsonData struct {
	Username  string `json:"username"`
	Room      string `json:"room"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type Client struct {
	username string
	hub      *ConnHub
	conn     *websocket.Conn
	send     chan []byte
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
			fmt.Println("Error reading JSON", err)
		}
		fmt.Printf("Got response %#v\n", message)

		messageJson, _ := json.Marshal(message)
		// queue messge for writing
		c.hub.broadcast <- messageJson
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// channel has been closed by the hub
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// coalesce pending messages into one message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
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

func (ws WebSocketServer) InitWebSocket() {
	ws.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
}

func (ws WebSocketServer) WsHandler(hub *ConnHub, w http.ResponseWriter, r *http.Request) {

}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{}
}
