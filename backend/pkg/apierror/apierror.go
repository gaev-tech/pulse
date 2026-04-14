package apierror

import (
	"encoding/json"
	"net/http"
)

type detail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type response struct {
	Error detail `json:"error"`
}

func Write(writer http.ResponseWriter, status int, code, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(response{
		Error: detail{Code: code, Message: message},
	})
}

func BadRequest(writer http.ResponseWriter, message string) {
	Write(writer, http.StatusBadRequest, "VALIDATION_ERROR", message)
}

func Unauthorized(writer http.ResponseWriter) {
	Write(writer, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
}

func Forbidden(writer http.ResponseWriter) {
	Write(writer, http.StatusForbidden, "PERMISSION_DENIED", "permission denied")
}

func NotFound(writer http.ResponseWriter) {
	Write(writer, http.StatusNotFound, "NOT_FOUND", "not found")
}

func Conflict(writer http.ResponseWriter, message string) {
	Write(writer, http.StatusConflict, "CONFLICT", message)
}

func QuotaExceeded(writer http.ResponseWriter, message string) {
	Write(writer, http.StatusPaymentRequired, "QUOTA_EXCEEDED", message)
}

func Internal(writer http.ResponseWriter) {
	Write(writer, http.StatusInternalServerError, "INTERNAL", "internal server error")
}
