package main

import (
	_ "base/src/infrastructure/database"
	httpserver "base/src/infrastructure/http_server"
)

type AppModule struct{}

func (appModule *AppModule) InitServer() {
	// HTTP SERVER
	httpServer := httpserver.NewHttpServer()
	httpServer.Init()
	httpServer.RegisterRoutes(Routes())
	httpServer.Listen()
}

func NewAppModule() *AppModule {
	return &AppModule{}
}
