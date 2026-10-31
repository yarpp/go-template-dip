package handler

import (
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	userService service.UserServiceInterface
}

func NewAuthHandler(userService service.UserServiceInterface) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req model.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		// Ошибки валидации всегда 400
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.userService.Register(r.Context(), &req)
	if err != nil {
		HandleServiceError(w, err)
		return
	}
	h.respondWithJSON(w, resp, http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req model.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.userService.Login(r.Context(), &req)
	if err != nil {
		HandleServiceError(w, err)
		return
	}
	h.respondWithJSON(w, resp, http.StatusOK)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Опционально — в роутере не используется
	h.respondWithError(w, "not implemented", http.StatusNotImplemented)
}

func (h *AuthHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *AuthHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	WriteError(w, message, statusCode)
}
