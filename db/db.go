package db

import (
	"database/sql"
	"log"
	"sync"

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

// --- In-Memory Хранилище ---
type MemoryStorage struct {
	mu       sync.RWMutex
	posts    map[string]*model.Post
	comments map[string][]*model.Comment
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		posts:    make(map[string]*model.Post),
		comments: make(map[string][]*model.Comment),
	}
}

func (m *MemoryStorage) GetPostByID(id string) (*model.Post, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	post, exists := m.posts[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	return post, nil
}

func (m *MemoryStorage) GetAllPosts() ([]*model.Post, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*model.Post
	for _, post := range m.posts {
		result = append(result, post)
	}
	return result, nil
}

func (m *MemoryStorage) CreatePost(post *model.Post) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.posts[post.ID] = post
	return nil
}

func (m *MemoryStorage) CreateComment(comment *model.Comment) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	postID := comment.PostID
	if _, exists := m.comments[postID]; !exists {
		m.comments[postID] = []*model.Comment{comment}
	} else {
		m.comments[postID] = append(m.comments[postID], comment)
	}
	return nil
}

func (m *MemoryStorage) GetCommentsByPostID(postID string, limit, offset int) ([]*model.Comment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	comments, exists := m.comments[postID]
	if !exists {
		return nil, sql.ErrNoRows
	}
	if offset >= len(comments) {
		return nil, nil
	}
	if limit == 0 || limit > len(comments)-offset {
		limit = len(comments) - offset
	}
	return comments[offset : offset+limit], nil
}

// PostgreSQL Хранилище
type PostgresStorage struct {
	DB *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{DB: db}, nil
}

func (p *PostgresStorage) GetPostByID(id string) (*model.Post, error) {
	query := "SELECT id, title, content, author, comments_allowed, created_at FROM posts WHERE id = $1"
	row := p.DB.QueryRow(query, id)

	var post model.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CommentsAllowed, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (p *PostgresStorage) GetAllPosts() ([]*model.Post, error) {
	rows, err := p.DB.Query("SELECT id, title, content, author, comments_allowed, created_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CommentsAllowed, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (p *PostgresStorage) CreatePost(post *model.Post) error {
	query := "INSERT INTO posts (id, title, content, author, comments_allowed, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := p.DB.Exec(query, post.ID, post.Title, post.Content, post.Author, post.CommentsAllowed, post.CreatedAt)
	return err
}

func (p *PostgresStorage) CreateComment(comment *model.Comment) error {
	query := "INSERT INTO comments (id, post_id, content, author, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := p.DB.Exec(query, comment.ID, comment.PostID, comment.Content, comment.Author, comment.CreatedAt)
	return err
}

func (p *PostgresStorage) GetCommentsByPostID(postID string, limit, offset int) ([]*model.Comment, error) {
	query := "SELECT id, post_id, content, author, created_at FROM comments WHERE post_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"
	rows, err := p.DB.Query(query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Content, &comment.Author, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

// Инициализация хранилища
func InitStorage(cfg *config.Config) {
	if cfg.StorageType == "postgres" {
		db, err := NewPostgresStorage(cfg.PostgresDSN)
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
