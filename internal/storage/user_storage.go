package storage

import (
	"fmt"
	"time"
)

func NewUserStore(path string) (*UserStore, error) {
	s := &UserStore{path: path}
	if err := loadJSON(path, &s.Users); err != nil {
		return nil, err
	}
	for _, u := range s.Users {
		if u.ID > s.next {
			s.next = u.ID
		}
	}
	return s, nil
}

func (s *UserStore) Create(username, email, passwordHash string) (*UserRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, u := range s.Users {
		if u.Email == email {
			return nil, fmt.Errorf("email already exists")
		}
		if u.Username == username {
			return nil, fmt.Errorf("username already exists")
		}
	}

	s.next++
	now := time.Now().UTC()
	u := &UserRecord{
		ID:        s.next,
		Username:  username,
		Email:     email,
		Password:  passwordHash,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.Users = append(s.Users, u)
	if err := atomicWrite(s.path, s.Users); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserStore) GetByID(id int) (*UserRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.Users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

func (s *UserStore) GetByEmail(email string) (*UserRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.Users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, ErrNotFound
}
