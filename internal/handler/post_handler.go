package handler

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/logger"
	"advanced-blog-management-system/internal/middleware"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type PostHandler struct {
	postService service.PostServiceInterface
	eventLogger *logger.EventLogger
}

func NewPostHandler(postService service.PostServiceInterface, eventLogger *logger.EventLogger) *PostHandler {
	return &PostHandler{postService: postService, eventLogger: eventLogger}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, apperrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}
	var req model.PostCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	post, err := h.postService.Create(r.Context(), userID, &req)
	if err != nil {
		HandleServiceError(w, err)
		return
	}
	h.respondWithJSON(w, post, http.StatusCreated)
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		WriteError(w, "invalid post id", http.StatusBadRequest)
		return
	}
	requestorID, _ := middleware.GetUserIDFromContext(r.Context())
	post, err := h.postService.GetByID(r.Context(), id, requestorID)
	if err != nil {
		HandleServiceError(w, err)
		return
	}
	h.respondWithJSON(w, post, http.StatusOK)
}

func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	posts, total, err := h.postService.GetAll(r.Context(), limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}
	h.respondWithJSON(w, map[string]interface{}{
		"posts":  posts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}, http.StatusOK)
}

func (h *PostHandler) GetByAuthor(w http.ResponseWriter, r *http.Request) {
	authorStr := chi.URLParam(r, "authorID")
	authorID, err := strconv.Atoi(authorStr)
	if err != nil || authorID <= 0 {
		WriteError(w, "invalid author id", http.StatusBadRequest)
		return
	}
	limit, offset := parsePagination(r)
	posts, total, err := h.postService.GetByAuthor(r.Context(), authorID, limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}
	h.respondWithJSON(w, map[string]interface{}{
		"posts":  posts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}, http.StatusOK)
}

func (h *PostHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *PostHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	WriteError(w, message, statusCode)
}

func parsePagination(r *http.Request) (int, int) {
	limit := 10
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			if n > 100 {
				n = 100
			}
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return limit, offset
}
