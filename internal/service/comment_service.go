package service

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/logger"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/storage"
	"context"
	"fmt"
	"strings"
)

type CommentService struct {
	commentStore *storage.CommentStore
	postStore    *storage.PostStore
	eventLogger  *logger.EventLogger
}

func NewCommentService(commentStore *storage.CommentStore, postStore *storage.PostStore, eventLogger *logger.EventLogger) *CommentService {
	return &CommentService{commentStore: commentStore, postStore: postStore, eventLogger: eventLogger}
}

func (s *CommentService) Create(ctx context.Context, userID, postID int, content string) (*model.Comment, error) {
	if postID <= 0 {
		return nil, apperrors.ErrInvalidPostID
	}
	if _, err := s.postStore.GetByID(postID); err != nil {
		return nil, apperrors.ErrPostNotFound
	}
	content = strings.TrimSpace(content)
	if content == "" || len([]rune(content)) > 1000 {
		return nil, apperrors.ErrValidation
	}

	c, err := s.commentStore.Create(content, postID, userID)
	if err != nil {
		return nil, err
	}
	if s.eventLogger != nil {
		s.eventLogger.LogEvent(fmt.Sprintf("user %d created comment %d", userID, c.ID))
	}
	return &model.Comment{
		ID:        c.ID,
		Content:   c.Content,
		PostID:    c.PostID,
		AuthorID:  c.AuthorID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}, nil
}

func (s *CommentService) GetByPost(ctx context.Context, postID, limit, offset int) ([]*model.Comment, int, error) {
	if postID <= 0 {
		return nil, 0, apperrors.ErrInvalidPostID
	}
	if _, err := s.postStore.GetByID(postID); err != nil {
		return nil, 0, apperrors.ErrPostNotFound
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	records, total := s.commentStore.GetByPost(postID, limit, offset)
	comments := make([]*model.Comment, 0, len(records))
	for _, c := range records {
		comments = append(comments, &model.Comment{
			ID:        c.ID,
			Content:   c.Content,
			PostID:    c.PostID,
			AuthorID:  c.AuthorID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}
	return comments, total, nil
}
