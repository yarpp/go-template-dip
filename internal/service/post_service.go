package service

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/logger"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/storage"
	"context"
	"fmt"
)

type PostService struct {
	postStore   *storage.PostStore
	userStore   *storage.UserStore
	eventLogger *logger.EventLogger
}

func NewPostService(postStore *storage.PostStore, userStore *storage.UserStore, eventLogger *logger.EventLogger) *PostService {
	return &PostService{postStore: postStore, userStore: userStore, eventLogger: eventLogger}
}

func (s *PostService) Create(ctx context.Context, userID int, req *model.PostCreateRequest) (*model.Post, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	p, err := s.postStore.Create(req.Title, req.Content, userID)
	if err != nil {
		return nil, err
	}
	// отправка события в канал логирования
	if s.eventLogger != nil {
		s.eventLogger.LogEvent(fmt.Sprintf("user %d created post %d", userID, p.ID))
	}
	return &model.Post{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		AuthorID:  p.AuthorID,
		CreatedAt: p.CreatedAt,
	}, nil
}

func (s *PostService) GetByID(ctx context.Context, id int, requestorID int) (*model.Post, error) {
	p, err := s.postStore.GetByID(id)
	if err != nil {
		return nil, apperrors.ErrPostNotFound
	}
	return &model.Post{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		AuthorID:  p.AuthorID,
		CreatedAt: p.CreatedAt,
	}, nil
}

func (s *PostService) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	records, total := s.postStore.GetAll(limit, offset)
	posts := make([]*model.Post, 0, len(records))
	for _, p := range records {
		posts = append(posts, &model.Post{
			ID:        p.ID,
			Title:     p.Title,
			Content:   p.Content,
			AuthorID:  p.AuthorID,
			CreatedAt: p.CreatedAt,
		})
	}
	return posts, total, nil
}

func (s *PostService) GetByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	records, total := s.postStore.GetByAuthor(authorID, limit, offset)
	posts := make([]*model.Post, 0, len(records))
	for _, p := range records {
		posts = append(posts, &model.Post{
			ID:        p.ID,
			Title:     p.Title,
			Content:   p.Content,
			AuthorID:  p.AuthorID,
			CreatedAt: p.CreatedAt,
		})
	}
	return posts, total, nil
}

// PostServiceInterface методы уже совпадают с сигнатурами шаблона.
