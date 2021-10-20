package main

import (
	"net/http"

	messages_usecases "base/src/applications/usecases/messages"
	"base/src/infrastructure/logger"
	repositories "base/src/infrastructure/repositories/messages_repository"
	"base/src/interfaces/http/controllers/messages"
	"base/src/interfaces/http/middlewares"
	controllerbase "base/src/shared"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Routes() *http.Handler {

	router := mux.NewRouter()
	subrouter := router.Host("localhost:8000").PathPrefix("/api/v1").Subrouter()

	/* Logger */
	logger := logger.NewLogger()

	/* Middlewares */
	loggingMiddleware := middlewares.NewLoggingMiddleware(logger)
	subrouter.Use(loggingMiddleware.LoggingMiddleware)

	/* Response Factory */
	httpResponseFactory := controllerbase.NewHttpResponseFactory()

	/* Messages routes */
	messagesRepository := repositories.NewMessagesRepository()
	messagesUseCases := messages_usecases.NewMessageUseCase(messagesRepository)
	messageController := messages.NewMessagesController(*httpResponseFactory, messagesUseCases)

	subrouter.HandleFunc("/messages", messageController.Post).Methods(http.MethodPost)
	subrouter.HandleFunc("/messages", messageController.Index).Methods(http.MethodGet)
	subrouter.HandleFunc("/messages/getbyroom/{room}", messageController.GetByRoom).Methods(http.MethodGet)
	subrouter.HandleFunc("/messages/{id:[0-9]+}", messageController.Show).Methods(http.MethodGet)
	subrouter.HandleFunc("/messages/{id:[0-9]+}", messageController.Put).Methods(http.MethodPut)
	subrouter.HandleFunc("/messages/{id:[0-9]+}", messageController.Delete).Methods(http.MethodDelete)

	handler := cors.Default().Handler(subrouter)

	return &handler

}
