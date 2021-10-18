package main

import (
	"net/http"

	usecases "base/src/applications/usecases/messages"
	repositories "base/src/infrastructure/repositories/messages_repository"
	"base/src/interfaces/http/controllers/messages"
	controllerbase "base/src/shared"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Routes() *http.Handler {

	router := mux.NewRouter()
	subrouter := router.Host("localhost:8000").PathPrefix("/api/v1").Subrouter()

	/* Response Factory */
	httpResponseFactory := controllerbase.NewHttpResponseFactory()

	/* Messages routes */
	messagesRepository := repositories.NewMessagesRepository()
	messagesUseCases := usecases.NewMessageUseCase(messagesRepository)
	messageController := messages.NewMessagesController(*httpResponseFactory, messagesUseCases)

	subrouter.HandleFunc("/messages", messageController.Post).Methods(http.MethodPost)
	subrouter.HandleFunc("/messages", messageController.Index).Methods(http.MethodGet)
	subrouter.HandleFunc("/messages/{id:[0-9]+}", messageController.Show).Methods(http.MethodGet)
	subrouter.HandleFunc("/messages/{id:[0-9]+}", messageController.Put).Methods(http.MethodPut)
	subrouter.HandleFunc("/messages/{id:[0-9]+}", messageController.Delete).Methods(http.MethodDelete)

	handler := cors.Default().Handler(subrouter)

	return &handler

}
