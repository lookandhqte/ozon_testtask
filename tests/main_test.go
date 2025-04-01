package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"ozon_test/graph"
	"ozon_test/storage"

	"github.com/stretchr/testify/assert"
)

func setupTestDB() {
	storage.DB = storage.NewMemoryStorage()
}

func TestMain(m *testing.M) {
	setupTestDB()
	os.Exit(m.Run())
}

// Тест создания поста
func TestCreatePost(t *testing.T) {
	setupTestDB()

	resolver := &graph.Resolver{}
	ctx := context.Background()

	newPost, err := resolver.Mutation().CreatePost(ctx, "Тестовое название", "Контент", "Автор", true)
	assert.NoError(t, err)
	assert.NotNil(t, newPost)
	assert.Equal(t, "Тестовое название", newPost.Title)
	assert.Equal(t, "Контент", newPost.Content)
	assert.Equal(t, "Автор", newPost.Author)
	assert.True(t, newPost.CommentsAllowed)
}

// Тест добавления комментария
func TestAddComment(t *testing.T) {
	setupTestDB()

	resolver := &graph.Resolver{}
	ctx := context.Background()

	newPost, err := resolver.Mutation().CreatePost(ctx, "Тест", "Контент", "Автор", true)
	assert.NoError(t, err)

	newComment, err := resolver.Mutation().AddComment(ctx, newPost.ID, nil, "Автор коммента", "Тестовый коммент")
	assert.NoError(t, err)
	assert.NotNil(t, newComment)
	assert.Equal(t, newPost.ID, newComment.PostID)
}

// Тест получения комментариев к посту
func TestGetPostWithComments(t *testing.T) {
	setupTestDB()

	resolver := &graph.Resolver{}
	ctx := context.Background()

	newPost, _ := resolver.Mutation().CreatePost(ctx, "Тест", "Контент", "Автор", true)
	newComment, _ := resolver.Mutation().AddComment(ctx, newPost.ID, nil, "Автор", "Комментарий")

	comments, err := resolver.Query().Comments(ctx, newPost.ID, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, newComment.ID, comments[0].ID)
}

// Тест создания поста без комментариев
func TestCreatePostWithoutComments(t *testing.T) {
	setupTestDB()

	resolver := &graph.Resolver{}
	ctx := context.Background()

	newPost, err := resolver.Mutation().CreatePost(ctx, "Пост без комментов", "Контент", "Автор", false)
	assert.NoError(t, err)
	assert.False(t, newPost.CommentsAllowed)
}

// Тест запрета добавления комментариев
func TestAddCommentToPostWithDisabledComments(t *testing.T) {
	setupTestDB()

	resolver := &graph.Resolver{}
	ctx := context.Background()

	newPost, _ := resolver.Mutation().CreatePost(ctx, "Без комментов", "Контент", "Автор", false)
	newComment, err := resolver.Mutation().AddComment(ctx, newPost.ID, nil, "Автор", "Коммент")

	assert.Error(t, err)
	assert.Nil(t, newComment)
	assert.Equal(t, "комментарии к этому посту запрещены", err.Error())
}

// Интеграционный тест полного процесса
func TestFullProcess(t *testing.T) {
	setupTestDB()

	resolver := &graph.Resolver{}
	ctx := context.Background()

	newPost, _ := resolver.Mutation().CreatePost(ctx, "Тестовый пост", "Контент", "Автор", true)
	newComment, _ := resolver.Mutation().AddComment(ctx, newPost.ID, nil, "Тест", "Комментарий")

	time.Sleep(100 * time.Millisecond)

	retrievedPost, err := resolver.Query().Post(ctx, newPost.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedPost)
	assert.Equal(t, newPost.ID, retrievedPost.ID)

	assert.NotNil(t, retrievedPost.Comments)
	assert.Len(t, retrievedPost.Comments, 1)
	assert.Equal(t, newComment.ID, retrievedPost.Comments[0].ID)
}

// Тест пагинации
func TestPaginationComments(t *testing.T) {
	setupTestDB()

	resolver := &graph.Resolver{}
	ctx := context.Background()

	newPost, _ := resolver.Mutation().CreatePost(ctx, "Pagination", "Content", "Author", true)

	for i := 1; i <= 5; i++ {
		_, _ = resolver.Mutation().AddComment(ctx, newPost.ID, nil, fmt.Sprintf("User %d", i), fmt.Sprintf("Comment %d", i))
	}

	time.Sleep(100 * time.Millisecond)

	firstPage, _ := resolver.Query().Comments(ctx, newPost.ID, 2, 0)
	assert.Len(t, firstPage, 2)

	secondPage, _ := resolver.Query().Comments(ctx, newPost.ID, 2, 2)
	assert.Len(t, secondPage, 2)
	assert.NotEqual(t, firstPage[0].ID, secondPage[0].ID)

	thirdPage, _ := resolver.Query().Comments(ctx, newPost.ID, 2, 4)
	assert.Len(t, thirdPage, 1)

	emptyPage, _ := resolver.Query().Comments(ctx, newPost.ID, 2, 6)
	assert.Empty(t, emptyPage)
}

// Тест вложенных комментариев
func TestNestedComments(t *testing.T) {
	setupTestDB()

	resolver := &graph.Resolver{}
	ctx := context.Background()

	post, _ := resolver.Mutation().CreatePost(ctx, "Вложенность", "Контент", "Автор", true)

	root, _ := resolver.Mutation().AddComment(ctx, post.ID, nil, "User1", "Root")
	child1, _ := resolver.Mutation().AddComment(ctx, post.ID, &root.ID, "User2", "Child 1")
	child2, _ := resolver.Mutation().AddComment(ctx, post.ID, &child1.ID, "User3", "Child 2")

	comments, err := resolver.Query().Comments(ctx, post.ID, 10, 0)
	assert.NoError(t, err)
	assert.NotNil(t, child2)
	assert.Len(t, comments, 3)
}

// Тест подписки на комментарии
func TestCommentSubscription(t *testing.T) {
	setupTestDB()

	resolver := &graph.Resolver{}
	ctx := context.Background()

	post, _ := resolver.Mutation().CreatePost(ctx, "Подписка", "Контент", "Автор", true)

	commentChan, err := resolver.Subscription().CommentAdded(ctx, post.ID)
	assert.NoError(t, err)

	go func() {
		time.Sleep(1 * time.Second)
		_, _ = resolver.Mutation().AddComment(ctx, post.ID, nil, "User", "Новый комментарий")
	}()

	select {
	case comment := <-commentChan:
		assert.NotNil(t, comment)
		assert.Equal(t, "Новый комментарий", comment.Content)
	case <-time.After(3 * time.Second):
		t.Fatal("Не получили комментарий через подписку")
	}
}
