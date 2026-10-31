package model

import (
	"time"
	//"github.com/go-playground/validator/v10"
)

// User представляет модель пользователя в системе
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Post представляет модель поста в блоге
type Post struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Comment представляет модель комментария к посту
type Comment struct {
	ID        int       `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	PostID    int       `json:"post_id" db:"post_id"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserCreateRequest представляет запрос на создание пользователя
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserLoginRequest представляет запрос на вход пользователя
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// PostCreateRequest представляет запрос на создание поста
type PostCreateRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

// CommentCreateRequest представляет запрос на создание комментария
type CommentCreateRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
	PostID  int    `json:"post_id" validate:"required,gt=0"`
}

// UserResponse - структура для ответа с данными пользователя (без пароля)
type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenResponse - структура для ответа с JWT токеном
type TokenResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// PostResponse - структура для ответа с данными поста
type PostResponse struct {
	ID        int          `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
}

// CommentResponse - структура для ответа с данными комментария
type CommentResponse struct {
	ID        int          `json:"id"`
	Content   string       `json:"content"`
	PostID    int          `json:"post_id"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// TODO: Реализовать ToResponse()
// Преобразовать User в UserResponse (скопировать поля, исключая Password)
func (u *User) ToResponse() UserResponse {
	// TODO: реализовать
	return UserResponse{}
}

// TODO: Реализовать Validate() для UserCreateRequest
// Использовать validator.New().Struct(r)
func (r *UserCreateRequest) Validate() error {
	// TODO: реализовать
	return nil
}

// TODO: Реализовать Validate() для UserLoginRequest
// Использовать validator.New().Struct(r)
func (r *UserLoginRequest) Validate() error {
	// TODO: реализовать
	return nil
}

// TODO: Реализовать Validate() для PostCreateRequest
// Использовать validator.New().Struct(r)
func (r *PostCreateRequest) Validate() error {
	// TODO: реализовать
	return nil
}

// TODO: Реализовать Validate() для CommentCreateRequest
// Использовать validator.New().Struct(r)
func (r *CommentCreateRequest) Validate() error {
	// TODO: реализовать
	return nil
}
