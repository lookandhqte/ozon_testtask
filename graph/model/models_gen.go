// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Comment struct {
	ID        string  `json:"id"`
	PostID    string  `json:"postId"`
	ParentID  *string `json:"parentId,omitempty"`
	Author    string  `json:"author"`
	Content   string  `json:"content"`
	CreatedAt string  `json:"createdAt"`
}

type Mutation struct {
}

type Post struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	Author          string     `json:"author"`
	CommentsAllowed bool       `json:"commentsAllowed"`
	CreatedAt       string     `json:"createdAt"`
	Comments        []*Comment `json:"comments,omitempty"`
}

type Query struct {
}

type Subscription struct {
}