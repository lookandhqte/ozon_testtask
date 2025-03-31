package storage

import (
	"log"
	"ozon_test/config"
	"ozon_test/graph/model"

	_ "github.com/lib/pq"
)

// Интерфейс хранилища (общий для PostgreSQL и in-memory)
type Storage interface {
	GetPostByID(id string) (*model.Post, error)
	GetAllPosts() ([]*model.Post, error)
	CreatePost(post *model.Post) error
	CreateComment(comment *model.Comment) error
	GetCommentsByPostID(postID string, limit, offset int) ([]*model.Comment, error)
}

// Глобальная переменная для хранения выбранного хранилища
var DB Storage

// Инициализация хранилища
func InitStorage(cfg *config.Config) {
	if cfg.StorageType == "postgres" {
		db, err := NewPostgresStorage(cfg.DSN)
		if err != nil {
			log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
		}
		DB = db
		log.Println("Используется PostgreSQL для хранения данных")
	} else {
		DB = NewMemoryStorage()
		log.Println("Используется In-Memory хранилище")
	}
}
