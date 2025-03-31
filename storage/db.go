package storage

import (
	"database/sql"
	"ozon_test/graph/model"
)

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
