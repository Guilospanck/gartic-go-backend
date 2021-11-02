package messages

import (
	"base/src/business/dtos"
	usecases_interfaces "base/src/business/usecases"
	controllerbase "base/src/shared"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type IMessagesController interface {
	Post(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	GetByRoom(w http.ResponseWriter, r *http.Request)
	Index(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	DeleteByRoom(w http.ResponseWriter, r *http.Request)
}

type messagesController struct {
	httpResponseFactory controllerbase.HttpResponseFactory
	usecases            usecases_interfaces.IMessagesUseCases
}

func (mc messagesController) Post(w http.ResponseWriter, r *http.Request) {
	createUserDTO := dtos.CreateMessageDTO{}

	err := json.NewDecoder(r.Body).Decode(&createUserDTO)
	if err != nil {
		res := mc.httpResponseFactory.BadRequest(err.Error(), nil)
		w.WriteHeader(res.StatusCode)
		json.NewEncoder(w).Encode(res)
		return
	}

	result, err := mc.usecases.CreateMessage(createUserDTO)
	if err != nil {
		res := mc.httpResponseFactory.BadRequest(err.Error(), nil)
		w.WriteHeader(res.StatusCode)
		json.NewEncoder(w).Encode(res)
		return
	}

	response := mc.httpResponseFactory.Created(result, nil)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func (mc messagesController) Show(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("show"))
}

func (mc messagesController) GetByRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	room := vars["room"]

	result, err := mc.usecases.GetMessagesByRoom(room)
	if err != nil {
		res := mc.httpResponseFactory.InternalServerError(err.Error(), nil)
		w.WriteHeader(res.StatusCode)
		json.NewEncoder(w).Encode(res)
		return
	}

	response := mc.httpResponseFactory.Ok(result, nil)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func (mc messagesController) Index(w http.ResponseWriter, r *http.Request) {
	result, err := mc.usecases.GetAllMessages()
	if err != nil {
		res := mc.httpResponseFactory.InternalServerError(err.Error(), nil)
		w.WriteHeader(res.StatusCode)
		json.NewEncoder(w).Encode(res)
		return
	}

	response := mc.httpResponseFactory.Ok(result, nil)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func (mc messagesController) Put(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("put"))
}

func (mc messagesController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("delete"))
}

func (mc messagesController) DeleteByRoom(w http.ResponseWriter, r *http.Request) {
	room := r.Header.Get("room")

	err := mc.usecases.DeleteAllMessagesFromRoom(room)
	if err != nil {
		res := mc.httpResponseFactory.InternalServerError(err.Error(), nil)
		w.WriteHeader(res.StatusCode)
		return
	}

	w.Write([]byte("OK"))
}

func NewMessagesController(httpResponseFactory controllerbase.HttpResponseFactory, messagesUsecase usecases_interfaces.IMessagesUseCases) IMessagesController {
	return &messagesController{
		httpResponseFactory: httpResponseFactory,
		usecases:            messagesUsecase,
	}
}
