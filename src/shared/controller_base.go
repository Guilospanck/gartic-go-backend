package controllerbase

import (
	httpreqres "base/src/infrastructure"
	"net/http"
)

type HttpResponseFactory struct{}

type ErrorMessage struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (*HttpResponseFactory) Ok(body interface{}, headers http.Header) httpreqres.HttpResponse {
	return httpreqres.HttpResponse{
		StatusCode: 200,
		Body:       body,
		Headers:    headers,
	}
}

func (*HttpResponseFactory) Created(body interface{}, headers http.Header) httpreqres.HttpResponse {
	return httpreqres.HttpResponse{
		StatusCode: 201,
		Body:       body,
		Headers:    headers,
	}
}

func (*HttpResponseFactory) NoContent(headers http.Header) httpreqres.HttpResponse {
	return httpreqres.HttpResponse{
		StatusCode: 204,
		Headers:    headers,
	}
}

func (*HttpResponseFactory) BadRequest(msg string, headers http.Header) httpreqres.HttpResponse {
	return httpreqres.HttpResponse{
		StatusCode: 400,
		Body: ErrorMessage{
			StatusCode: 400,
			Message:    msg,
		},
		Headers: headers,
	}
}

func (*HttpResponseFactory) Unauthorized(msg string, headers http.Header) httpreqres.HttpResponse {
	return httpreqres.HttpResponse{
		StatusCode: 401,
		Body: ErrorMessage{
			StatusCode: 401,
			Message:    msg,
		},
		Headers: headers,
	}
}

func (*HttpResponseFactory) Forbidden(msg string, headers http.Header) httpreqres.HttpResponse {
	return httpreqres.HttpResponse{
		StatusCode: 403,
		Body: ErrorMessage{
			StatusCode: 403,
			Message:    msg,
		},
		Headers: headers,
	}
}

func (*HttpResponseFactory) NotFound(msg string, headers http.Header) httpreqres.HttpResponse {
	return httpreqres.HttpResponse{
		StatusCode: 404,
		Body: ErrorMessage{
			StatusCode: 404,
			Message:    msg,
		},
		Headers: headers,
	}
}

func (*HttpResponseFactory) Conflict(msg string, headers http.Header) httpreqres.HttpResponse {
	return httpreqres.HttpResponse{
		StatusCode: 409,
		Body: ErrorMessage{
			StatusCode: 409,
			Message:    msg,
		},
		Headers: headers,
	}
}

func (*HttpResponseFactory) InternalServerError(msg string, headers http.Header) httpreqres.HttpResponse {
	return httpreqres.HttpResponse{
		StatusCode: 500,
		Body: ErrorMessage{
			StatusCode: 500,
			Message:    msg,
		},
		Headers: headers,
	}
}

func NewHttpResponseFactory() *HttpResponseFactory {
	return &HttpResponseFactory{}
}
