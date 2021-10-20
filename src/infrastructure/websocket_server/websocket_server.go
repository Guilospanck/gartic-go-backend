package websocketserver

import (
	httpserver "base/src/infrastructure/http_server"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type WebSocketServer struct {
	httpServerPort  uint
	httpserver      *httpserver.HttpServer
	httpHandlerFunc *http.Handler
}

func (ws WebSocketServer) InitHttpServer() {
	httpServer := httpserver.NewHttpServer(ws.httpServerPort)
	ws.httpserver = httpServer

	ws.httpserver.Init()
}

func (ws WebSocketServer) InitWebSocket() {
	router := mux.NewRouter()
	router.HandleFunc("/ws", serveWs)

	handler := cors.Default().Handler(router)
	ws.httpHandlerFunc = &handler
}

func (ws WebSocketServer) Listen() {
	ws.httpserver.RegisterRoutes(ws.httpHandlerFunc)
	ws.httpserver.Listen()
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func NewWebSocketServer(httpServerPort ...uint) *WebSocketServer {
	httpDefaultPort := 5555
	if len(httpServerPort) == 0 {
		httpServerPort = append(httpServerPort, uint(httpDefaultPort))
	}

	return &WebSocketServer{
		httpServerPort: httpServerPort[0],
	}
}
