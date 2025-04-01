# Post & Comment System (GraphQL + Go + PostgreSQL)

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-4169E1?logo=postgresql)
![GraphQL](https://img.shields.io/badge/GraphQL-E10098?logo=graphql)

Микросервис для управления постами и комментариями с поддержкой GraphQL API, реализованный на Go с возможностью выбора хранилища данных (PostgreSQL или In-Memory).

## 📌 Основные возможности

### Посты

- Создание постов с настраиваемой политикой комментариев
- Просмотр списка постов с пагинацией
- Детализация отдельного поста

### Комментарии

- Иерархическая система комментариев (вложенность)
- Автоматическая валидация длины (до 2000 символов)
- Пагинация комментариев
- Режим "только для чтения" для постов

### Реальное время

- GraphQL Subscriptions для мгновенных обновлений
- Уведомления о новых комментариях

### Инфраструктура

- Два режима хранения данных:
  - **PostgreSQL** - для production
  - **In-Memory** - для разработки и тестирования
- Полная контейнеризация (Docker)
- Интеграционные и unit-тесты

## 🛠 Технологический стек

| Компонент       | Технологии                        |
| --------------- | --------------------------------- |
| Бэкенд          | Go 1.21+                          |
| GraphQL         | gqlgen                            |
| Базы данных     | PostgreSQL 17 / In-Memory storage |
| Контейнеризация | Docker + Docker Compose           |
| Тестирование    | `go test` + testify/assert        |

## 🚀 Быстрый старт

### Запуск через Docker

Флаг --profile отвечает за выбор хранилища. Доступен в двух вариантах: postgres, memory

1. $ docker-compose --profile postgres up --build | **PostgreSQL режим:**
2. $ docker-compose --profile memory up --build | **In-memory режим:**

### Запуск через консоль

1. $ STORAGE_TYPE=postgres
   $ POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/ozon_test?sslmode=disable
   $ go run main.go http://localhost:8080
2. $ STORAGE_TYPE=memory
   $ go run main.go http://localhost:8080

### Основные запросы

Получить запросы
#query { posts { id title content author commentsAllowed createdAt } }

Создать новый пост
#mutation { createPost( title: "" content: "" author: "" commentsAllowed: true ) { id title content } }

Получить конкретный пост
#query { post(id: "") { id title content author commentsAllowed createdAt } }

Добавить комментарий к посту
#mutation { addComment( postId: "" parentId: null author: "" content: "" ) { id content author createdAt } }

Подписка на новые комментарии
#subscription { commentAdded(postId: "") { id content author createdAt } }

Получить все комментарии к посту
->variables: { "postId": "", "limit": 10, "offset": 0 }
#query GetCommentsByPost($postId: ID!, $limit: Int = 100, $offset: Int = 0) { comments(postID: $postId, limit: $limit, offset: $offset) { id postId parentId author content createdAt } }

### Тестирование

$ go test -v .\tests\...
