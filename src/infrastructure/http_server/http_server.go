package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type HttpServer struct {
	server *http.Server
	port   int
}

func (hs *HttpServer) Init() {
	hs.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", hs.port),
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

func NewHttpServer(port ...int) *HttpServer {
	/* A way of adding optional parameters. See: https://stackoverflow.com/a/19813113/9782182 */
	defaultPort := 8000
	if len(port) == 0 {
		port = append(port, int(defaultPort))
	}

	return &HttpServer{
		port: port[0],
	}
}
