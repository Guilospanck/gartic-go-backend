package websocketserver

import (
	"base/src/business/dtos"
	usecases_interfaces "base/src/business/usecases"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Client struct {
	Username string           `json:"username"`
	Room     string           `json:"room"`
	Hub      *ConnHub         `json:"hub"`
	Conn     *websocket.Conn  `json:"conn"`
	Send     chan interface{} `json:"send"`
}

type JsonData struct {
	Username          string `json:"username"`
	Room              string `json:"room"`
	Message           string `json:"message"`
	Date              string `json:"date"`
	Close             bool   `json:"close"`
	CanvasCoordinates string `json:"canvasCoordinates"`
	CanvasConfigs     string `json:"canvasConfigs"`
	DidPlayerWin      bool   `json:"didPlayerWin"`
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
	Drawing   string   `json:"drawing"`
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
	roomWithParticipants := client.getParticipantsFromRoom(client.Room)
	client.Hub.broadcastToParticipantsInRoom <- roomWithParticipants
}

func (client *Client) getRandomParticipant() ParticipantsTurn {
	clients := client.Hub.clients[client.Room]
	randomIndex := rand.Intn(len(clients))

	randomClient := clients[randomIndex]
	timestamp := time.Now().Unix()

	drawing := client.getRandomDrawing()

	response := ParticipantsTurn{
		Username:  randomClient.Username,
		Room:      randomClient.Room,
		Hub:       randomClient.Hub,
		Timestamp: timestamp,
		Drawing:   drawing,
	}

	return response
}

func (client *Client) getRandomDrawing() string {
	drawings := []string{"flower", "potato", "sun", "car", "human", "lighter", "santa"}
	randomIndex := rand.Intn(len(drawings))
	return drawings[randomIndex]
}

func (c *Client) ReadPump(messageUsecase usecases_interfaces.IMessagesUseCases) {
	// schedule client to be disconnected
	defer func() {
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

		// verify if is a winner
		if message.DidPlayerWin {
			// TODO: do something
		}

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
		if err != nil {
			fmt.Errorf("error trying to save message to the database")
		}

		// queue messge for writing
		c.Hub.send <- message
	}
}

func (c *Client) WritePump(messageUsecase usecases_interfaces.IMessagesUseCases, drawerUsecase usecases_interfaces.IDrawersUseCases) {
	ticker := time.NewTicker(pingPeriod)
	sendParticipantsTurn := time.NewTicker(drawersTimer)

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
				fmt.Errorf("Write pump not ok")
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			messageJson, _ := json.Marshal(message)
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
			if c.Room == "waitingroomgarticlikeapp" {
				break
			}

			drawer, err := drawerUsecase.GetDrawerByRoom(c.Room)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Errorf("error record not found")

				client := c.getRandomParticipant()

				newDrawer := dtos.CreateDrawerDTO{
					Username: client.Username,
					Room:     client.Room,
				}
				drawerUsecase.CreateDrawer(newDrawer)

				c.Hub.broadcastParticipantTurn <- client

				break
			}

			if err != nil {
				fmt.Errorf("error trying to get drawer by room")
			}

			drawerTimestamp := drawer.CreatedAt
			now := time.Now()
			timePassed := now.Sub(drawerTimestamp)

			if timePassed.Seconds() < float64(drawersTimer.Seconds()-1) {
				break
			}

			numberOfParticipants := len(c.Hub.clients[c.Room])
			client := c.getRandomParticipant()
			if numberOfParticipants > 1 {
				for client.Username == drawer.Username {
					client = c.getRandomParticipant()
				}
			}

			drawerUsecase.DeleteAllDrawersFromRoom(c.Room)
			newDrawer := dtos.CreateDrawerDTO{
				Username: client.Username,
				Room:     client.Room,
			}
			drawerUsecase.CreateDrawer(newDrawer)

			c.Hub.broadcastParticipantTurn <- client
		}
	}
}
