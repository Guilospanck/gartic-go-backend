package httpserver

import (
	"bytes"
	"net/http"
)

type HttpResponseFactory struct {
	w    http.ResponseWriter
	buf  bytes.Buffer
	code int
}

type ErrorMessage struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (rw *HttpResponseFactory) Header() http.Header {
	return rw.w.Header()
}

func (rw *HttpResponseFactory) Write(b []byte) (int, error) {
	switch responseType := string(b); responseType {
	case "Ok":
		return rw.w.Write([]byte("Test"))
	default:
		return rw.w.Write([]byte("Default"))
	}
}

func (rw *HttpResponseFactory) WriteHeader(statusCode int) {
	rw.code = statusCode
}

func Ok(body interface{}, headers http.Header) HttpResponse {
	return HttpResponse{
		StatusCode: 200,
		Body:       body,
		Headers:    headers,
	}
}

func Created(body interface{}, headers http.Header) HttpResponse {
	return HttpResponse{
		StatusCode: 201,
		Body:       body,
		Headers:    headers,
	}
}

func NoContent(headers http.Header) HttpResponse {
	return HttpResponse{
		StatusCode: 204,
		Headers:    headers,
	}
}

func BadRequest(msg string, headers http.Header) HttpResponse {
	return HttpResponse{
		StatusCode: 400,
		Body: ErrorMessage{
			StatusCode: 400,
			Message:    msg,
		},
		Headers: headers,
	}
}

func Unauthorized(msg string, headers http.Header) HttpResponse {
	return HttpResponse{
		StatusCode: 401,
		Body: ErrorMessage{
			StatusCode: 401,
			Message:    msg,
		},
		Headers: headers,
	}
}

func Forbidden(msg string, headers http.Header) HttpResponse {
	return HttpResponse{
		StatusCode: 403,
		Body: ErrorMessage{
			StatusCode: 403,
			Message:    msg,
		},
		Headers: headers,
	}
}

func NotFound(msg string, headers http.Header) HttpResponse {
	return HttpResponse{
		StatusCode: 404,
		Body: ErrorMessage{
			StatusCode: 404,
			Message:    msg,
		},
		Headers: headers,
	}
}

func Conflict(msg string, headers http.Header) HttpResponse {
	return HttpResponse{
		StatusCode: 409,
		Body: ErrorMessage{
			StatusCode: 409,
			Message:    msg,
		},
		Headers: headers,
	}
}

func InternalServerError(msg string, headers http.Header) HttpResponse {
	return HttpResponse{
		StatusCode: 500,
		Body: ErrorMessage{
			StatusCode: 500,
			Message:    msg,
		},
		Headers: headers,
	}
}

func NewHttpResponseFactory(w http.ResponseWriter) *HttpResponseFactory {
	return &HttpResponseFactory{
		w: w,
	}
}
