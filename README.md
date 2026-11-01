#  Blog Management System - Template Project

Это шаблон дипломного проекта на Go для разработки REST API блог-платформы с регистрацией пользователей, управлением постами и комментариями.

## 📋 Содержание

- [О проекте](#о-проекте)
- [Структура проекта](#структура-проекта)
- [Технологический стек](#технологический-стек)
- [Быстрый старт](#быстрый-старт)
- [Разработка](#разработка)
- [API эндпоинты](#api-эндпоинты)
- [Примеры запросов](#примеры-запросов)

## О проекте

** Blog Management System** - это дипломный проект для студентов программы "Go-разработчик с нуля".

Проект демонстрирует:
- ✅ Основы HTTP сервера и REST API
- ✅ Работу с JSON и структурами Go
- ✅ Регистрацию и авторизацию пользователей
- ✅ Управление постами и комментариями
- ✅ Хеширование паролей (bcrypt)
- ✅ JWT токены
- ✅ Отложенное логирование через горутину и канал
- ✅ Структурирование проекта по директориям
- ✅ Контейнеризацию в Docker

## Структура проекта

```
template_project/
├── api/
│   └── main.go                 # Точка входа приложения
├── internal/
│   ├── handler/                # HTTP обработчики (handlers)
│   │   ├── auth_handler.go
│   │   ├── post_handler.go
│   │   └── comment_handler.go
│   ├── middleware/             # HTTP middleware (auth, logging)
│   │   ├── auth.go
│   │   └── logging.go
│   ├── model/                  # Модели данных (структуры)
│   │   └── models.go
│   ├── storage/                # Работа с JSON файлами
│   │   ├── user_storage.go
│   │   ├── post_storage.go
│   │   └── comment_storage.go
│   └── logger/                 # Event logger (логирование в файл)
│       └── event_logger.go
|
├── pkg/
│   └── auth/                   # Утилиты аутентификации
│       ├── jwt.go              # JWT токены
│       └── password.go         # Хеширование паролей (bcrypt)
├── data/                       # JSON файлы с данными
│   ├── users.json
│   ├── posts.json
│   └── comments.json
├── logs.txt                    # Логи событий (создается при запуске)
├── .env.example                # Пример конфигурации
├── docker-compose.yml          # Docker Compose
├── Dockerfile                  # Docker образ
└── go.mod                      # Зависимости проекта
```

## Технологический стек

- **Язык:** Go 1.21+
- **Стандартная библиотека:** net/http, encoding/json
- **Аутентификация:** JWT (golang-jwt/jwt)
- **Хеширование:** bcrypt (golang.org/x/crypto)
- **Конфигурация:** godotenv
- **Логирование:** горутины и каналы
- **Контейнеризация:** Docker, Docker Compose

## Быстрый старт

### Предварительные требования

- Go 1.21 или выше
- Docker и Docker Compose (опционально)

### 1. Подготовка окружения

```bash
# Клонировать репозиторий
git clone <repo-url>
cd template_project

# Скопировать конфигурацию
cp .env.example .env

# Установить Go зависимости
go mod download
```

### 2. Разработка и реализация

После реализации всех компонентов по TODO комментариям:

```bash
# Запустить приложение локально
go run api/main.go

# Приложение будет доступно на http://localhost:8080
```

### 3. Запуск в Docker (опционально)

```bash
# Собрать и запустить в Docker
docker-compose up --build

# Остановить
docker-compose down
```

## Разработка

### Что уже готово ✅

- Структура проекта и директории
- Модели данных (User, Post, Comment)
- Функции хеширования паролей (HashPassword, CheckPassword)
- Основные HTTP обработчики с TODO
- docker-compose.yml для контейнеризации

### Что нужно реализовать ❌

Проект содержит TODO комментарии, которые указывают что реализовать:

#### 1. **Storage (работа с JSON файлами)** - `internal/storage/`

Создайте файлы для работы с данными:
- `user_storage.go` - сохранение/загрузка пользователей
- `post_storage.go` - сохранение/загрузка постов
- `comment_storage.go` - сохранение/загрузка комментариев

Функции для реализации:
- Загрузка данных из JSON при старте
- Сохранение новых объектов в JSON
- Поиск объектов по ID
- Проверка существования по email/username
- Синхронизация (написать в файл сразу)

#### 2. **JWT токены** - `pkg/auth/jwt.go`

- `GenerateToken(userID int, email string) (string, error)` - создание токена
- `ValidateToken(tokenString string) (int, error)` - проверка токена, возврат userID

#### 3. **Middleware** - `internal/middleware/`

- `AuthMiddleware` - проверка JWT токена в заголовке Authorization
- `LoggingMiddleware` - логирование HTTP запросов

#### 4. **Обработчики** - `internal/handler/`

Реализуйте методы обработчиков:

**AuthHandler:**
- `Register(w http.ResponseWriter, r *http.Request)` - POST /register
- `Login(w http http.ResponseWriter, r *http.Request)` - POST /login

**PostHandler:**
- `Create(w http.ResponseWriter, r *http.Request)` - POST /posts
- `GetAll(w http.ResponseWriter, r *http.Request)` - GET /posts
- `GetByID(w http.ResponseWriter, r *http.Request)` - GET /posts/{id}

**CommentHandler:**
- `Create(w http.ResponseWriter, r *http.Request)` - POST /posts/{id}/comments
- `GetByPost(w http.ResponseWriter, r *http.Request)` - GET /posts/{id}/comments

#### 5. **Event Logger** - `internal/logger/event_logger.go`

Логирование создания постов и комментариев:
- `NewEventLogger(filePath string) *EventLogger` - инициализация
- `LogEvent(event string)` - отправка события в канал
- `Start()` - запуск worker горутины
- `Stop()` - остановка логера при завершении приложения
- `worker()` - горутина, которая записывает логи в файл с задержкой

#### 6. **Главная функция** - `cmd/api/main.go`

Инициализируйте:
- Загрузку переменных окружения (.env)
- JWT менеджер
- Storage (пользователей, постов, комментариев)
- Event Logger
- Обработчики
- Middleware
- HTTP маршруты
- HTTP сервер
- Graceful shutdown (обработка Ctrl+C)

## API эндпоинты

### Публичные эндпоинты

```
GET    /api/health                     # Проверка здоровья API
POST   /api/register                   # Регистрация пользователя
POST   /api/login                      # Вход пользователя
GET    /api/posts                      # Получить все посты
GET    /api/posts/{id}                 # Получить пост по ID
GET    /api/posts/{id}/comments        # Получить комментарии к посту
```

### Защищенные эндпоинты (требуют Authorization: Bearer TOKEN)

```
POST   /api/posts                      # Создать пост
POST   /api/posts/{id}/comments        # Добавить комментарий к посту
```

## Примеры запросов

### Health Check
```bash
curl http://localhost:8080/api/health
```

**Ответ:**
```json
{
  "status": "ok"
}
```

### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Ответ (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### Вход пользователя
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Ответ (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### Создание поста (требуется токен)
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first post"
  }'
```

**Ответ (201 Created):**
```json
{
  "id": 1,
  "title": "My First Post",
  "content": "This is the content of my first post",
  "author_id": 1,
  "created_at": "2024-01-15T10:35:00Z"
}
```

### Получение всех постов
```bash
curl http://localhost:8080/api/posts
```

**Ответ (200 OK):**
```json
[
  {
    "id": 1,
    "title": "My First Post",
    "content": "This is the content of my first post",
    "author_id": 1,
    "created_at": "2024-01-15T10:35:00Z"
  }
]
```

### Получение конкретного поста
```bash
curl http://localhost:8080/api/posts/1
```

**Ответ (200 OK):**
```json
{
  "id": 1,
  "title": "My First Post",
  "content": "This is the content of my first post",
  "author_id": 1,
  "created_at": "2024-01-15T10:35:00Z"
}
```

### Добавление комментария (требуется токен)
```bash
curl -X POST http://localhost:8080/api/posts/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "content": "Great post!"
  }'
```

**Ответ (201 Created):**
```json
{
  "id": 1,
  "content": "Great post!",
  "post_id": 1,
  "author_id": 2,
  "created_at": "2024-01-15T10:40:00Z"
}
```

### Получение комментариев к посту
```bash
curl http://localhost:8080/api/posts/1/comments
```

**Ответ (200 OK):**
```json
[
  {
    "id": 1,
    "content": "Great post!",
    "post_id": 1,
    "author_id": 2,
    "created_at": "2024-01-15T10:40:00Z"
  }
]
```

## Конфигурация

Переменные окружения задаются в файле `.env`:

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRY_HOURS=24

# Data storage
DATA_DIR=./data
LOGS_FILE=./logs.txt

# Environment
ENV=development
```

## Логирование событий

Приложение логирует события создания постов и комментариев в файл `logs.txt`:

```
[2024-01-15 10:35:45] user 1 created post 1
[2024-01-15 10:40:20] user 2 created comment 1
```

Логирование реализовано с использованием:
- **Канала** (channel) для отправки событий
- **Горутины** (goroutine) для асинхронной записи в файл
- **Задержки** (sleep) для демонстрации отложенной обработки

## Архитектура приложения

```
┌─────────────────┐
│  HTTP Requests  │
└────────┬────────┘
         │
┌────────▼────────────────────────┐
│ Middleware (Auth, Logging)      │
└────────┬────────────────────────┘
         │
┌────────▼─────────────┐
│ Handlers (HTTP API)  │ ← Парсинг JSON, валидация
└────────┬─────────────┘
         │
┌────────▼─────────────┐
│ Storage (JSON)       │ ← Сохранение/загрузка данных
└────────┬─────────────┘
         │
┌────────▼──────────────┐
│ JSON файлы            │ ← users.json, posts.json, comments.json
└───────────────────────┘
```

## Ключевые концепции

### Структуры (Structs)

Используйте теги для JSON сериализации:
```go
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`  // Не включать в JSON
    CreatedAt time.Time `json:"created_at"`
}
```

### JWT токены

Токены генерируются при регистрации/входе и проверяются в middleware:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Горутины и каналы

Event Logger использует канал для отправки событий и горутину для их обработки:
```go
eventsChan := make(chan string, 100)
go worker(eventsChan)  // Горутина записывает логи
eventsChan <- "user 1 created post 1"  // Отправка события
```

### Обработка ошибок

Всегда проверяйте ошибки и возвращайте правильные HTTP коды:
```go
if err != nil {
    if err == ErrUserNotFound {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
}
```

## Полезные команды

```bash
# Скачать зависимости
go mod download

# Запустить приложение
go run cmd/api/main.go

# Собрать приложение
go build -o api ./api/main.go
./api

# Запустить в Docker
docker-compose up --build

# Остановить Docker сервисы
docker-compose down

# Просмотреть логи приложения
tail -f logs.txt

# Проверить JSON файлы данных
cat data/users.json
cat data/posts.json
cat data/comments.json
```

## Требования к сдаче

Перед отправкой убедитесь, что:

- ✅ Все 6 публичных эндпоинтов работают
- ✅ Оба защищенных эндпоинта работают (требуют токен)
- ✅ Регистрация создает уникального пользователя
- ✅ Авторизация выдает JWT токен
- ✅ Посты и комментарии сохраняются в JSON файлы
- ✅ Логирование работает (события в logs.txt с задержкой)
- ✅ Приложение запускается через `go run` и Docker
- ✅ Коды ошибок правильные (400, 401, 404, 500 и т.д.)
- ✅ README актуален

## Полезные ссылки

- [Go Tour](https://tour.golang.org/) - интерактивное введение в Go
- [HTTP пакет в Go](https://pkg.go.dev/net/http) - официальная документация
- [JSON в Go](https://pkg.go.dev/encoding/json) - работа с JSON
- [JWT-go документация](https://github.com/golang-jwt/jwt) - JWT токены
- [bcrypt документация](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - хеширование паролей
