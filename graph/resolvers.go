package graph

//go:generate go run github.com/99designs/gqlgen generate
// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"errors"
	"time"

	"ozon_test/db"
	"ozon_test/graph/model"

	"github.com/google/uuid"
)

type Resolver struct{}

// Создание нового поста
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, author string, commentsAllowed bool) (*model.Post, error) {
	post := &model.Post{
		ID:              uuid.New().String(),
		Title:           title,
		Content:         content,
		Author:          author,
		CommentsAllowed: commentsAllowed,
		CreatedAt:       time.Now().Format(time.RFC3339),
	}

	err := db.DB.CreatePost(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// Добавление комментария
func (r *mutationResolver) AddComment(ctx context.Context, postID string, parentID *string, author, content string) (*model.Comment, error) {
	post, err := db.DB.GetPostByID(postID)

	if err != nil || post == nil {
		return nil, errors.New("пост не найден")
	}

	if !post.CommentsAllowed {
		return nil, errors.New("комментарии к этому посту запрещены")
	}
	if len(content) > 2000 {
		return nil, errors.New("Комментарий слишком длинный (максимум 2000 символов)")
	}

	comment := &model.Comment{
		ID:        uuid.New().String(),
		PostID:    postID,
		ParentID:  parentID,
		Author:    author,
		Content:   content,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	err = db.DB.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// Получение всех постов
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	return db.DB.GetAllPosts()
}

// Получение поста по ID
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	post, err := db.DB.GetPostByID(id)
	if err != nil || post == nil {
		return nil, errors.New("пост не найден")
	}

	// Загрузка комментариев к посту
	comments, err := db.DB.GetCommentsByPostID(id, 10, 0)
	if err != nil {
		return nil, err
	}
	post.Comments = comments // Добавляем комментарии к объекту поста

	return post, nil
}

// Получение комментариев к посту с поддержкой пагинации
func (r *queryResolver) Comments(ctx context.Context, postID string, limit int, offset int) ([]*model.Comment, error) {
	return db.DB.GetCommentsByPostID(postID, limit, offset)
}

// Поддержка подписки на новые комментарии (GraphQL Subscriptions)
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	commentChan := make(chan *model.Comment, 1)

	go func() {
		for {
			// Ожидаем новый комментарий в БД
			time.Sleep(1 * time.Second)
			comments, _ := db.DB.GetCommentsByPostID(postID, 1, 0) // Получаем последний комментарий
			if len(comments) > 0 {
				commentChan <- comments[0] // Отправляем реальный комментарий
			}
		}
	}()

	return commentChan, nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
