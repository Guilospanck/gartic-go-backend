package main

import (
	messages_usecases "base/src/applications/usecases/messages"
	_ "base/src/infrastructure/database"
	httpserver "base/src/infrastructure/http_server"
	repositories "base/src/infrastructure/repositories/messages_repository"
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
	go hub.Run()

	// Websocket Server (and its own http server)
	messagesRepository := repositories.NewMessagesRepository()
	messagesUseCases := messages_usecases.NewMessageUseCase(messagesRepository)

	webSocketServer := websocketserver.NewWebSocketServer(messagesUseCases)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		webSocketServer.WsHandler(hub, w, r)
	})
	log.Fatal(http.ListenAndServe(":5555", nil))

}

func NewAppModule() *AppModule {
	return &AppModule{}
}
