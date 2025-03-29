package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Post represents a blog post
type Post struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	Author          string    `json:"author"`
	CommentsEnabled bool      `json:"commentsEnabled"`
	CreatedAt       time.Time `json:"createdAt"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"postId"`
	ParentID  *string   `json:"parentId"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// Resolver handles GraphQL queries
type Resolver struct {
	posts      []*Post
	postsMutex sync.Mutex

	comments      []*Comment
	commentsMutex sync.Mutex

	commentObservers map[string]map[string]chan *Comment
	observersMutex   sync.Mutex
}

// NewResolver creates a new resolver instance
func NewResolver() *Resolver {
	return &Resolver{
		commentObservers: make(map[string]map[string]chan *Comment),
	}
}

// Posts returns all posts
func (r *Resolver) Posts(ctx context.Context) ([]*Post, error) {
	r.postsMutex.Lock()
	defer r.postsMutex.Unlock()

	return r.posts, nil
}

// Post returns a single post by ID
func (r *Resolver) Post(ctx context.Context, args struct{ ID string }) (*Post, error) {
	r.postsMutex.Lock()
	defer r.postsMutex.Unlock()

	for _, post := range r.posts {
		if post.ID == args.ID {
			return post, nil
		}
	}

	return nil, fmt.Errorf("post not found")
}

// CreatePost creates a new post
func (r *Resolver) CreatePost(ctx context.Context, args struct {
	Title           string
	Content         string
	Author          string
	CommentsEnabled bool
}) (*Post, error) {
	r.postsMutex.Lock()
	defer r.postsMutex.Unlock()

	post := &Post{
		ID:              fmt.Sprintf("%d", len(r.posts)+1),
		Title:           args.Title,
		Content:         args.Content,
		Author:          args.Author,
		CommentsEnabled: args.CommentsEnabled,
		CreatedAt:       time.Now(),
	}

	r.posts = append(r.posts, post)
	return post, nil
}

// UpdatePostCommentsEnabled updates the comments enabled status of a post
func (r *Resolver) UpdatePostCommentsEnabled(ctx context.Context, args struct {
	ID              string
	CommentsEnabled bool
}) (*Post, error) {
	r.postsMutex.Lock()
	defer r.postsMutex.Unlock()

	for _, post := range r.posts {
		if post.ID == args.ID {
			post.CommentsEnabled = args.CommentsEnabled
			return post, nil
		}
	}

	return nil, fmt.Errorf("post not found")
}

// Comments returns comments for a post with pagination
func (r *Resolver) Comments(ctx context.Context, args struct {
	PostID string
	Limit  *int
	Offset *int
}) ([]*Comment, error) {
	r.commentsMutex.Lock()
	defer r.commentsMutex.Unlock()

	var result []*Comment
	limit := 10
	offset := 0

	if args.Limit != nil {
		limit = *args.Limit
	}
	if args.Offset != nil {
		offset = *args.Offset
	}

	// Filter comments by post ID and parent ID (nil for top-level comments)
	for _, comment := range r.comments {
		if comment.PostID == args.PostID && comment.ParentID == nil {
			result = append(result, comment)
		}
	}

	// Apply pagination
	if offset >= len(result) {
		return []*Comment{}, nil
	}

	end := offset + limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end], nil
}

// Replies returns replies to a comment
func (r *Resolver) Replies(ctx context.Context, obj *Comment) ([]*Comment, error) {
	r.commentsMutex.Lock()
	defer r.commentsMutex.Unlock()

	var replies []*Comment
	for _, comment := range r.comments {
		if comment.ParentID != nil && *comment.ParentID == obj.ID {
			replies = append(replies, comment)
		}
	}

	return replies, nil
}

// CreateComment creates a new comment
func (r *Resolver) CreateComment(ctx context.Context, args struct {
	PostID   string
	ParentID *string
	Author   string
	Content  string
}) (*Comment, error) {
	if len(args.Content) > 2000 {
		return nil, errors.New("comment content exceeds maximum length of 2000 characters")
	}

	r.postsMutex.Lock()
	var post *Post
	for _, p := range r.posts {
		if p.ID == args.PostID {
			post = p
			break
		}
	}
	r.postsMutex.Unlock()

	if post == nil {
		return nil, errors.New("post not found")
	}

	if !post.CommentsEnabled {
		return nil, errors.New("comments are disabled for this post")
	}

	r.commentsMutex.Lock()
	defer r.commentsMutex.Unlock()

	comment := &Comment{
		ID:        fmt.Sprintf("%d", len(r.comments)+1),
		PostID:    args.PostID,
		ParentID:  args.ParentID,
		Author:    args.Author,
		Content:   args.Content,
		CreatedAt: time.Now(),
	}

	r.comments = append(r.comments, comment)

	// Notify subscribers
	r.observersMutex.Lock()
	if observers, ok := r.commentObservers[args.PostID]; ok {
		for _, ch := range observers {
			ch <- comment
		}
	}
	r.observersMutex.Unlock()

	return comment, nil
}

// CommentAdded subscription resolver
func (r *Resolver) CommentAdded(ctx context.Context, args struct{ PostID string }) (<-chan *Comment, error) {
	ch := make(chan *Comment, 1)

	r.observersMutex.Lock()
	if _, ok := r.commentObservers[args.PostID]; !ok {
		r.commentObservers[args.PostID] = make(map[string]chan *Comment)
	}
	id := fmt.Sprintf("%p", &ch)
	r.commentObservers[args.PostID][id] = ch
	r.observersMutex.Unlock()

	go func() {
		<-ctx.Done()
		r.observersMutex.Lock()
		delete(r.commentObservers[args.PostID], id)
		r.observersMutex.Unlock()
	}()

	return ch, nil
}
