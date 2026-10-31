package handler

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/logger"
	"advanced-blog-management-system/internal/middleware"
	"advanced-blog-management-system/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type CommentHandler struct {
	commentService service.CommentServiceInterface
	eventLogger    *logger.EventLogger
}

func NewCommentHandler(commentService service.CommentServiceInterface, eventLogger *logger.EventLogger) *CommentHandler {
	return &CommentHandler{commentService: commentService, eventLogger: eventLogger}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		WriteError(w, apperrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}
	postID, err := strconv.Atoi(chi.URLParam(r, "postId"))
	if err != nil || postID <= 0 {
		WriteError(w, "invalid post id", http.StatusBadRequest)
		return
	}
	var body struct {
		Content string `json:"content" validate:"required,min=1,max=1000"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if err := validator.New().Struct(body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	comment, err := h.commentService.Create(r.Context(), userID, postID, body.Content)
	if err != nil {
		HandleServiceError(w, err)
		return
	}
	h.respondWithJSON(w, comment, http.StatusCreated)
}

func (h *CommentHandler) GetByPost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(chi.URLParam(r, "postId"))
	if err != nil || postID <= 0 {
		WriteError(w, "invalid post id", http.StatusBadRequest)
		return
	}
	limit, offset := parsePagination(r)
	comments, total, err := h.commentService.GetByPost(r.Context(), postID, limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}
	h.respondWithJSON(w, map[string]interface{}{
		"comments": comments,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	}, http.StatusOK)
}

func (h *CommentHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *CommentHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	WriteError(w, message, statusCode)
}
