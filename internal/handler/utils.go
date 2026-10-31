package handler

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func HandleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrUserAlreadyExists):
		WriteError(w, err.Error(), http.StatusConflict)
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		WriteError(w, err.Error(), http.StatusUnauthorized)
	case errors.Is(err, apperrors.ErrUserNotFound):
		WriteError(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, apperrors.ErrPostNotFound):
		WriteError(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, apperrors.ErrCommentNotFound):
		WriteError(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, apperrors.ErrForbidden):
		WriteError(w, err.Error(), http.StatusForbidden)
	case errors.Is(err, apperrors.ErrUnauthorized):
		WriteError(w, err.Error(), http.StatusUnauthorized)
	case errors.Is(err, apperrors.ErrInvalidPostID):
		WriteError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, apperrors.ErrValidation):
		WriteError(w, err.Error(), http.StatusBadRequest)
	default:
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteError(w, "internal server error", http.StatusInternalServerError)
	}
}
