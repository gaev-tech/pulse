package apierror

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Write(writer http.ResponseWriter, status int, code, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(Error{Code: code, Message: message})
}

func BadRequest(writer http.ResponseWriter, message string) {
	Write(writer, http.StatusBadRequest, "BAD_REQUEST", message)
}

func Unauthorized(writer http.ResponseWriter) {
	Write(writer, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
}

func Forbidden(writer http.ResponseWriter) {
	Write(writer, http.StatusForbidden, "FORBIDDEN", "forbidden")
}

func NotFound(writer http.ResponseWriter) {
	Write(writer, http.StatusNotFound, "NOT_FOUND", "not found")
}

func Internal(writer http.ResponseWriter) {
	Write(writer, http.StatusInternalServerError, "INTERNAL", "internal server error")
}
