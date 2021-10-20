package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type HttpServer struct {
	server *http.Server
	port   uint
}

func (hs *HttpServer) Init() {
	hs.server = &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", hs.port),
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

func NewHttpServer(port ...uint) *HttpServer {
	/* A way of adding optional parameters. See: https://stackoverflow.com/a/19813113/9782182 */
	defaultPort := 8000
	if len(port) == 0 {
		port = append(port, uint(defaultPort))
	}

	return &HttpServer{
		port: port[0],
	}
}
