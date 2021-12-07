package websocketserver

import (
	"net/http"
	"time"

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
	// maxMessageSize = 512

	// time to change drawer's turn
	drawersTimer = 20 * time.Second
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
	messagesUseCases usecases_interfaces.IMessagesUseCases
	drawerUseCases   usecases_interfaces.IDrawersUseCases
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

	go client.WritePump(ws.messagesUseCases, ws.drawerUseCases)
	go client.ReadPump(ws.messagesUseCases)
}

func NewWebSocketServer(messagesUsecase usecases_interfaces.IMessagesUseCases,
	drawersUsecase usecases_interfaces.IDrawersUseCases) IWebSocketServer {
	return &webSocketServer{
		messagesUseCases: messagesUsecase,
		drawerUseCases:   drawersUsecase,
	}
}
