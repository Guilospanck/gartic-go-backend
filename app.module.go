package main

import (
	drawers_usecases "base/src/applications/usecases/drawers"
	messages_usecases "base/src/applications/usecases/messages"
	_ "base/src/infrastructure/database"
	httpserver "base/src/infrastructure/http_server"
	drawers_repositories "base/src/infrastructure/repositories/drawers_repository"
	messages_repositories "base/src/infrastructure/repositories/messages_repository"
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
	messagesRepository := messages_repositories.NewMessagesRepository()
	messagesUseCases := messages_usecases.NewMessageUseCase(messagesRepository)
	drawersRepository := drawers_repositories.NewDrawersRepository()
	drawerUseCases := drawers_usecases.NewDrawersUseCase(drawersRepository)

	webSocketServer := websocketserver.NewWebSocketServer(messagesUseCases, drawerUseCases)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		webSocketServer.WsHandler(hub, w, r)
	})
	log.Fatal(http.ListenAndServe(":5555", nil))

}

func NewAppModule() *AppModule {
	return &AppModule{}
}
