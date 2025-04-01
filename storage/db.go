package storage

import (
	"database/sql"
	"ozon_test/graph/model"
)

// PostgresStorage реализует интерфейс хранилища с использованием PostgreSQL.
type PostgresStorage struct {
	DB *sql.DB
}

// NewPostgresStorage создает подключение к PostgreSQL по переданному DSN.
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{DB: db}, nil
}

// GetPostByID возвращает пост по его ID.
func (p *PostgresStorage) GetPostByID(id string) (*model.Post, error) {
	const query = `
		SELECT id, title, content, author, comments_allowed, created_at
		FROM posts
		WHERE id = $1
	`
	row := p.DB.QueryRow(query, id)

	var post model.Post
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Author,
		&post.CommentsAllowed,
		&post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// GetAllPosts возвращает все посты из базы данных.
func (p *PostgresStorage) GetAllPosts() ([]*model.Post, error) {
	const query = `
		SELECT id, title, content, author, comments_allowed, created_at
		FROM posts
	`
	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Author,
			&post.CommentsAllowed,
			&post.CreatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (p *PostgresStorage) UpdatePost(post *model.Post) error {
	query := `UPDATE posts SET 
        comments_allowed = $1 
        WHERE id = $2`
	_, err := p.DB.Exec(query, post.CommentsAllowed, post.ID)
	return err
}

// CreatePost сохраняет новый пост в базе данных.
func (p *PostgresStorage) CreatePost(post *model.Post) error {
	const query = `
		INSERT INTO posts (id, title, content, author, comments_allowed, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := p.DB.Exec(query,
		post.ID,
		post.Title,
		post.Content,
		post.Author,
		post.CommentsAllowed,
		post.CreatedAt,
	)
	return err
}

// CreateComment сохраняет комментарий к посту.
func (p *PostgresStorage) CreateComment(comment *model.Comment) error {
	const query = `
		INSERT INTO comments (id, post_id, content, author, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := p.DB.Exec(query,
		comment.ID,
		comment.PostID,
		comment.Content,
		comment.Author,
		comment.CreatedAt,
	)
	return err
}

// GetCommentsByPostID возвращает список комментариев к посту с пагинацией.
func (p *PostgresStorage) GetCommentsByPostID(postID string, limit, offset int) ([]*model.Comment, error) {
	const query = `
		SELECT id, post_id, content, author, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := p.DB.Query(query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Content,
			&comment.Author,
			&comment.CreatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}
