package storage

import "time"

func NewCommentStore(path string) (*CommentStore, error) {
	s := &CommentStore{path: path}
	if err := loadJSON(path, &s.Comments); err != nil {
		return nil, err
	}
	for _, c := range s.Comments {
		if c.ID > s.next {
			s.next = c.ID
		}
	}
	return s, nil
}

func (s *CommentStore) Create(content string, postID, authorID int) (*CommentRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	now := time.Now().UTC()
	c := &CommentRecord{
		ID:        s.next,
		Content:   content,
		PostID:    postID,
		AuthorID:  authorID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.Comments = append(s.Comments, c)
	if err := atomicWrite(s.path, s.Comments); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CommentStore) GetByPost(postID, limit, offset int) ([]*CommentRecord, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	all := make([]*CommentRecord, 0)
	for _, c := range s.Comments {
		if c.PostID == postID {
			all = append(all, c)
		}
	}
	total := len(all)
	if offset >= total {
		return []*CommentRecord{}, total
	}
	all = all[offset:]
	if limit > 0 && limit < len(all) {
		all = all[:limit]
	}
	return all, total
}
