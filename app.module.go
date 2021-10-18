package main

import (
	"base/src/infrastructure/database"
	_ "base/src/infrastructure/environments"
	httpserver "base/src/infrastructure/http_server"
)

type AppModule struct{}

func (appModule *AppModule) InitServer() {
	// HTTP SERVER
	httpServer := httpserver.NewHttpServer()
	httpServer.Init()
	httpServer.RegisterRoutes(Routes())
	httpServer.Listen()

	// CONNECT TO DATABASE
	database.ConnectToDatabase()
}

func NewAppModule() *AppModule {
	return &AppModule{}
}
