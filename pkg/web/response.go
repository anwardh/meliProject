package web

import (
	"net/http"
	"strconv"
)

// Definição da Estrutura
type Response struct { // O omitempty fará com que, dependendo do que houver (erro ou não) um campo será omitido
	Code  string      `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func NewResponse(code int, data interface{}, err string) Response {

	if code < http.StatusMultipleChoices { // Status 300
		return Response{strconv.FormatInt(int64(code), 10), data, ""} // Omitindo o Error
	}
	return Response{strconv.FormatInt(int64(code), 10), nil, err} // Omitindo o Data
}
