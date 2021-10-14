package main

import (
	"log"
	"net/http"
	"time"

	httpserver "base/src/infrastructure/http_server"
	"base/src/interfaces/http/controllers/messages"
	"base/src/interfaces/http/middlewares"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func RegisterRoutes() {

	router := mux.NewRouter()
	subrouter := router.Host("localhost:8000").Subrouter()

	/* Messages routes */
	messageController := messages.NewMessagesController(httpserver.HttpResponseFactory{})

	subrouter.HandleFunc("/api/v1/messages", middlewares.ResponseFactory(messageController.Post)).Methods(http.MethodPost)
	subrouter.HandleFunc("/api/v1/messages", messageController.Index).Methods(http.MethodGet)
	subrouter.HandleFunc("/api/v1/messages/:id", messageController.Show).Methods(http.MethodGet)
	subrouter.HandleFunc("/api/v1/messages/:id", messageController.Put).Methods(http.MethodPut)
	subrouter.HandleFunc("/api/v1/messages/:id", messageController.Delete).Methods(http.MethodDelete)

	handler := cors.Default().Handler(subrouter)

	// Server
	server := &http.Server{
		Handler:      handler,
		Addr:         "localhost:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())

}
