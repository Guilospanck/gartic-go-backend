package httpserver

import (
	"log"
	"net/http"
	"time"
)

type HttpServer struct {
	server *http.Server
}

func (hs *HttpServer) Init() {
	hs.server = &http.Server{
		Addr:         "localhost:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func (hs *HttpServer) RegisterRoutes(handler *http.Handler) {
	hs.server.Handler = *handler
}

func (hs *HttpServer) Listen() {
	log.Fatal(hs.server.ListenAndServe())
}

func NewHttpServer() *HttpServer {
	return &HttpServer{}
}
