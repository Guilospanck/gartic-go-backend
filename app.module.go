package main

import (
	_ "base/src/infrastructure/database"
	httpserver "base/src/infrastructure/http_server"
	websocketserver "base/src/infrastructure/websocket_server"
)

type AppModule struct{}

func (appModule *AppModule) InitServer() {
	// HTTP SERVER
	httpServer := httpserver.NewHttpServer()
	httpServer.Init()
	httpServer.RegisterRoutes(Routes())
	httpServer.Listen()

	// Websocket Server (and its own http server)
	webSocketServer := websocketserver.NewWebSocketServer()
	webSocketServer.InitHttpServer()
	webSocketServer.InitWebSocket()

}

func NewAppModule() *AppModule {
	return &AppModule{}
}
