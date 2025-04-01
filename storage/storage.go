package storage

import (
	"log"
	"ozon_test/config"
	"ozon_test/graph/model"

	_ "github.com/lib/pq" // импорт драйвера PostgreSQL
)

// Storage — интерфейс абстракции над типами хранилищ (PostgreSQL, Memory).
type Storage interface {
	GetPostByID(id string) (*model.Post, error)
	GetAllPosts() ([]*model.Post, error)
	CreatePost(post *model.Post) error
	CreateComment(comment *model.Comment) error
	GetCommentsByPostID(postID string, limit, offset int) ([]*model.Comment, error)
	UpdatePost(post *model.Post) error
}

// DB — глобальное хранилище, инициализируемое при старте приложения.
var DB Storage

// InitStorage инициализирует глобальное хранилище на основе конфигурации.
// Поддерживает два варианта: PostgreSQL и in-memory.
func InitStorage(cfg *config.Config) {
	switch cfg.StorageType {
	case "postgres":
		db, err := NewPostgresStorage(cfg.DSN)
		if err != nil {
			log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
		}
		DB = db
		log.Println("Используется PostgreSQL для хранения данных")

	default:
		DB = NewMemoryStorage()
		log.Println("Используется In-Memory хранилище")
	}
}
