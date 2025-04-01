package storage

import (
	"database/sql"
	"fmt"
	"ozon_test/graph/model"
	"sync"
)

// MemoryStorage реализует хранилище данных в оперативной памяти.
type MemoryStorage struct {
	mu       sync.RWMutex
	posts    map[string]*model.Post
	comments map[string][]*model.Comment
}

// NewMemoryStorage создает новое in-memory хранилище.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		posts:    make(map[string]*model.Post),
		comments: make(map[string][]*model.Comment),
	}
}
func (m *MemoryStorage) UpdatePost(post *model.Post) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.posts[post.ID]; !exists {
		return fmt.Errorf("пост не найден")
	}

	m.posts[post.ID] = post
	return nil
}

// GetPostByID возвращает пост по ID или ошибку, если не найден.
func (m *MemoryStorage) GetPostByID(id string) (*model.Post, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	post, exists := m.posts[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	return post, nil
}

// GetAllPosts возвращает все посты.
func (m *MemoryStorage) GetAllPosts() ([]*model.Post, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	posts := make([]*model.Post, 0, len(m.posts))
	for _, post := range m.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

// CreatePost добавляет новый пост.
func (m *MemoryStorage) CreatePost(post *model.Post) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.posts[post.ID] = post
	return nil
}

// CreateComment добавляет комментарий к посту.
func (m *MemoryStorage) CreateComment(comment *model.Comment) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	postID := comment.PostID
	m.comments[postID] = append(m.comments[postID], comment)
	return nil
}

// GetCommentsByPostID возвращает комментарии к посту с пагинацией.
func (m *MemoryStorage) GetCommentsByPostID(postID string, limit, offset int) ([]*model.Comment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	comments, exists := m.comments[postID]
	if !exists {
		return nil, sql.ErrNoRows
	}

	// Проверка границ
	if offset >= len(comments) {
		return nil, nil
	}
	end := offset + limit
	if limit == 0 || end > len(comments) {
		end = len(comments)
	}

	return comments[offset:end], nil
}
