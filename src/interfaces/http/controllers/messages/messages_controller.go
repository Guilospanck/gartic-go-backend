package messages

import (
	httpserver "base/src/infrastructure/http_server"
	"net/http"
)

type IMessagesController interface {
	Post(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	Index(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type messagesController struct {
	httpserver.HttpResponseFactory
}

func (mc messagesController) Post(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok"))
}

func (mc messagesController) Show(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("show"))
}

func (mc messagesController) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}

func (mc messagesController) Put(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("put"))
}

func (mc messagesController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete"))
}

func NewMessagesController(responseFactory httpserver.HttpResponseFactory) IMessagesController {
	return &messagesController{
		HttpResponseFactory: responseFactory,
	}
}
