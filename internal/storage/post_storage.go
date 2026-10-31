package storage

import "time"

func NewPostStore(path string) (*PostStore, error) {
	s := &PostStore{path: path}
	if err := loadJSON(path, &s.Posts); err != nil {
		return nil, err
	}
	for _, p := range s.Posts {
		if p.ID > s.next {
			s.next = p.ID
		}
	}
	return s, nil
}

func (s *PostStore) Create(title, content string, authorID int) (*PostRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	p := &PostRecord{
		ID:        s.next,
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: time.Now().UTC(),
	}
	s.Posts = append(s.Posts, p)
	if err := atomicWrite(s.path, s.Posts); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PostStore) GetByID(id int) (*PostRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.Posts {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

// GetAll возвращает копию постов (свежие сначала) с пагинацией
func (s *PostStore) GetAll(limit, offset int) ([]*PostRecord, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	total := len(s.Posts)
	result := make([]*PostRecord, 0)
	// обход в обратном порядке — новые сверху
	for i := len(s.Posts) - 1; i >= 0; i-- {
		result = append(result, s.Posts[i])
	}
	if offset >= len(result) {
		return []*PostRecord{}, total
	}
	result = result[offset:]
	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}
	return result, total
}

func (s *PostStore) GetByAuthor(authorID, limit, offset int) ([]*PostRecord, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	all := make([]*PostRecord, 0)
	for i := len(s.Posts) - 1; i >= 0; i-- {
		if s.Posts[i].AuthorID == authorID {
			all = append(all, s.Posts[i])
		}
	}
	total := len(all)
	if offset >= total {
		return []*PostRecord{}, total
	}
	all = all[offset:]
	if limit > 0 && limit < len(all) {
		all = all[:limit]
	}
	return all, total
}
