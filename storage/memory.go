package storage

import (
	"database/sql"
	"ozon_test/graph/model"
	"sync"
)

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
