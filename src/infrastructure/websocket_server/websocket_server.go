package websocketserver

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"base/src/business/dtos"

	usecases_interfaces "base/src/business/usecases"

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
	// needed to allow connections from any origin for :3333 -> :5555
	CheckOrigin: func(r *http.Request) bool { return true },
}

type IWebSocketServer interface {
	WsHandler(hub *ConnHub, w http.ResponseWriter, r *http.Request)
}

type webSocketServer struct {
	usecases usecases_interfaces.IMessagesUseCases
}

type JsonData struct {
	Username          string `json:"username"`
	Room              string `json:"room"`
	Message           string `json:"message"`
	Date              string `json:"date"`
	Close             bool   `json:"close"`
	CanvasCoordinates string `json:"canvasCoordinates"`
	CanvasConfigs     string `json:"canvasConfigs"`
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

type ParticipantsTurn struct {
	Username  string   `json:"username"`
	Room      string   `json:"room"`
	Hub       *ConnHub `json:"hub"`
	Timestamp int64    `json:"timestamp"`
}

func (client *Client) getParticipantsFromRoom(room string) RoomAndParticipants {
	participants := client.Hub.clients[room]

	participantsString := []string{}
	for _, v := range participants {
		participantsString = append(participantsString, v.Username)
	}

	roomsWithParticipants := RoomAndParticipants{
		Room:         room,
		Participants: participantsString,
	}

	return roomsWithParticipants
}

func (client *Client) getAllParticipantsFromRooms() []RoomAndParticipants {
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

	return roomsWithParticipants
}

func (client *Client) verifyHowManyParticipantsAreInTheRoom(room string) int {
	roomsWithParticipants := client.getAllParticipantsFromRooms()

	var numOfParticipants int

	for _, object := range roomsWithParticipants {
		if object.Room == room {
			numOfParticipants = len(object.Participants)
			break
		}
	}

	return numOfParticipants
}

func (client *Client) sendDataToWaitingRoom() {
	fmt.Println("Sending to waiting room...")

	roomsWithParticipants := client.getAllParticipantsFromRooms()

	client.Hub.sendToWaitingRoom <- roomsWithParticipants
}

func (client *Client) sendAllMessagesFromRoom(messageUsecase usecases_interfaces.IMessagesUseCases) {
	result, err := messageUsecase.GetMessagesByRoom(client.Room)
	if err != nil {
		return
	}

	client.Send <- result
}

func (client *Client) broadcastToParticipantsInRoom() {
	fmt.Printf("Broadcasting to room %s\n", client.Room)

	roomWithParticipants := client.getParticipantsFromRoom(client.Room)

	client.Hub.broadcastToParticipantsInRoom <- roomWithParticipants
}

func (client *Client) getRandomParticipant() ParticipantsTurn {
	clients := client.Hub.clients[client.Room]
	randomIndex := rand.Intn(len(clients))

	randomClient := clients[randomIndex]
	timestamp := time.Now().Unix()

	response := ParticipantsTurn{
		Username:  randomClient.Username,
		Room:      randomClient.Room,
		Hub:       randomClient.Hub,
		Timestamp: timestamp,
	}

	return response
}

func (c *Client) ReadPump(messageUsecase usecases_interfaces.IMessagesUseCases) {
	fmt.Println("listening...")

	// schedule client to be disconnected
	defer func() {
		fmt.Println("Read Pump: defer func")

		// Delete messages of the room from database if this is the last participant
		numOfParticipants := c.verifyHowManyParticipantsAreInTheRoom(c.Room)
		if numOfParticipants == 1 {
			messageUsecase.DeleteAllMessagesFromRoom(c.Room)
		}

		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	// init Client connection
	// c.Conn.SetReadLimit(maxMessageSize)
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
			return
		}

		// Save message to database
		messageDB := dtos.CreateMessageDTO{
			Username:          message.Username,
			Message:           message.Message,
			Room:              message.Room,
			Date:              message.Date,
			CanvasCoordinates: message.CanvasCoordinates,
			CanvasConfigs:     message.CanvasConfigs,
		}

		_, err := messageUsecase.CreateMessage(messageDB)
		if err == nil {
			fmt.Println("Saved to database")
		}

		// queue messge for writing
		c.Hub.send <- message
	}
}

func (c *Client) WritePump(messageUsecase usecases_interfaces.IMessagesUseCases) {
	fmt.Println("writing...")

	ticker := time.NewTicker(pingPeriod)
	sendParticipantsTurn := time.NewTicker(20 * time.Second)

	defer func() {
		ticker.Stop()
		c.sendDataToWaitingRoom()
		c.broadcastToParticipantsInRoom()
		c.Conn.Close()
	}()

	c.sendDataToWaitingRoom()
	if c.Room != "waitingroomgarticlikeapp" {
		c.broadcastToParticipantsInRoom()
		c.sendAllMessagesFromRoom(messageUsecase)
	}

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
			fmt.Println(messageJson)
			w.Write(messageJson)

			if err := w.Close(); err != nil {
				return
			}

		// send ping over websocket
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}

		// send participant's turn
		case <-sendParticipantsTurn.C:
			client := c.getRandomParticipant()
			client.Hub.broadcastParticipantTurn <- client

			// c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			// if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(clientJson)); err != nil {
			// 	return
			// }

		}
	}
}

func (ws webSocketServer) WsHandler(hub *ConnHub, w http.ResponseWriter, r *http.Request) {
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

	go client.WritePump(ws.usecases)
	go client.ReadPump(ws.usecases)
}

func NewWebSocketServer(messagesUsecase usecases_interfaces.IMessagesUseCases) IWebSocketServer {
	return &webSocketServer{
		usecases: messagesUsecase,
	}
}
