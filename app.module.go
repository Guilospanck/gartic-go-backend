package main

import (
	_ "base/src/infrastructure/database"
	httpserver "base/src/infrastructure/http_server"
	websocketserver "base/src/infrastructure/websocket_server"
	"log"
	"net/http"
)

type AppModule struct{}

func (appModule *AppModule) InitServer() {
	// HTTP SERVER
	httpServer := httpserver.NewHttpServer()
	httpServer.Init()
	httpServer.RegisterRoutes(Routes())
	go httpServer.Listen()

	// ========== WEB SOCKET ============
	// Start Hub
	hub := websocketserver.NewConnHub()
	hub.Run()

	// Websocket Server (and its own http server)
	webSocketServer := websocketserver.NewWebSocketServer()
	webSocketServer.InitWebSocket()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		webSocketServer.WsHandler(hub, w, r)
	})
	log.Fatal(http.ListenAndServe(":5555", nil))

}

func NewAppModule() *AppModule {
	return &AppModule{}
}
