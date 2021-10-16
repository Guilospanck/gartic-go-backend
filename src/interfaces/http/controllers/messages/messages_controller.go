package messages

import (
	controllerbase "base/src/shared"
	"encoding/json"
	"log"
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
	httpResponseFactory controllerbase.HttpResponseFactory
}

func (mc messagesController) Post(w http.ResponseWriter, r *http.Request) {

	body := r.Body
	headers := r.Header

	response := mc.httpResponseFactory.Created(body, headers)

	marshalledRes, err := json.Marshal(response)
	if err != nil {
		log.Fatal("Error")
	}

	w.WriteHeader(response.StatusCode)
	w.Write([]byte(marshalledRes))
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

func NewMessagesController(httpResponseFactory controllerbase.HttpResponseFactory) IMessagesController {
	return &messagesController{
		httpResponseFactory: httpResponseFactory,
	}
}
