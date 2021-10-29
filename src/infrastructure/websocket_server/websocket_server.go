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
	Close     bool   `json:"close"`
}

type Client struct {
	Username string           `json:"username"`
	Room     string           `json:"room"`
	Hub      *ConnHub         `json:"hub"`
	Conn     *websocket.Conn  `json:"conn"`
	Send     chan interface{} `json:"send"`
}

type RoomAndParticipants struct {
	Room         string   `json:"room"`
	Participants []string `json:"participants"`
}

func (client *Client) sendDataToWaitingRoom() {
	fmt.Println("Sending to waiting room...")

	roomsWithParticipants := []RoomAndParticipants{}

	for key, value := range client.Hub.clients {
		participants := []string{}
		for _, v := range value {
			participants = append(participants, v.Username)
		}

		temp := RoomAndParticipants{
			Room:         key,
			Participants: participants,
		}
		roomsWithParticipants = append(roomsWithParticipants, temp)
	}
	client.Hub.sendToWaitingRoom <- roomsWithParticipants
}

func (c *Client) ReadPump() {
	fmt.Println("listening...")

	// schedule client to be disconnected
	defer func() {
		fmt.Println("Read Pump: defer func")
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	// init Client connection
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(appData string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// handle connection read
	for {
		// read JSON data from connection
		message := JsonData{}
		if err := c.Conn.ReadJSON(&message); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		fmt.Printf("Got response %#v\n", message)

		// Verify if is to close the channel (will be closed by defer)
		if message.Close {
			break
		}

		// queue messge for writing
		c.Hub.send <- message
	}
}

func (c *Client) WritePump() {
	fmt.Println("writing...")

	if c.Room == "waitingroomgarticlikeapp" {
		c.sendDataToWaitingRoom()
	}

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
				// channel has been closed by the hub
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				fmt.Println("Write pump not ok")
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			messageJson, _ := json.Marshal(message)
			fmt.Println("Sent")
			w.Write(messageJson)

			// Add queued chat messages to the current websocket message.
			// n := len(c.Send)
			// for i := 0; i < n; i++ {
			// 	messageJson, _ := json.Marshal(<-c.Send)
			// 	w.Write(messageJson)
			// }

			if err := w.Close(); err != nil {
				return
			}

		// send ping over websocket
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
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
		Username: username,
		Room:     room,
		Conn:     conn,
		Hub:      hub,
		Send:     make(chan interface{}, 512),
	}

	client.Hub.register <- client
	gotoReadPump := make(chan int, 1)
	gotoSendDataToWaitingRoom := make(chan int, 1)

	go func() {
		go client.WritePump()
		gotoReadPump <- 1
	}()

	go func() {
		<-gotoReadPump
		go client.ReadPump()
		gotoSendDataToWaitingRoom <- 1
	}()

	go func() {
		<-gotoSendDataToWaitingRoom
		// go client.sendDataToWaitingRoom()
	}()

}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{}
}
