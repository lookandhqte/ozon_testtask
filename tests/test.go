package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"time"

	"ozon_test/graph"
	"ozon_test/storage"

	//"github.com/lookandhqte/ozon_test/graph/model"
	//"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() {
	storage.DB = storage.NewMemoryStorage()
}

func TestMain(m *testing.M) {
	setupTestDB() // Инициализация хранилища перед запуском тестов
	os.Exit(m.Run())
}

// Тест создания поста
func TestCreatePost(t *testing.T) {
	// Создаем in-memory хранилище перед тестом
	storage.DB = storage.NewMemoryStorage()

	resolver := &graph.Resolver{} // Используем указатель, чтобы передавать структуру по ссылке
	ctx := context.Background()

	newPost, err := resolver.Mutation().CreatePost(ctx, "Тестовое название", "Контент", "автор", true)

	assert.NoError(t, err, "Ошибка при создании поста")
	assert.NotNil(t, newPost, "Созданный пост не должен быть nil")
	assert.Equal(t, "Тестовое название", newPost.Title)
	assert.Equal(t, "Контент", newPost.Content)
	assert.Equal(t, "автор", newPost.Author)
	assert.Equal(t, true, newPost.CommentsAllowed)
}

// Тест добавления комментария
func TestAddComment(t *testing.T) {
	storage.DB = storage.NewMemoryStorage() // Инициализируем хранилище

	resolver := &graph.Resolver{}
	ctx := context.Background()

	// Создаем тестовый пост перед добавлением комментария
	newPost, err := resolver.Mutation().CreatePost(ctx, "Тестовое название", "Контент", "Автор", true)
	assert.NoError(t, err, "Ошибка при создании поста")
	assert.NotNil(t, newPost, "Созданный пост не должен быть nil")

	// Теперь создаем комментарий к только что созданному посту
	newComment, err := resolver.Mutation().AddComment(ctx, newPost.ID, nil, "Автор коммента", "Тестовый коммент")

	assert.NoError(t, err, "Ошибка при добавлении комментария")
	assert.NotNil(t, newComment, "Созданный комментарий не должен быть nil")
	assert.Equal(t, newPost.ID, newComment.PostID, "Комментарий должен быть привязан к правильному посту")
	assert.Equal(t, "Автор коммента", newComment.Author)
	assert.Equal(t, "Тестовый коммент", newComment.Content)
}

// Тест получения поста с комментариями
func TestGetPostWithComments(t *testing.T) {
	storage.DB = storage.NewMemoryStorage() // Инициализируем хранилище

	resolver := &graph.Resolver{}
	ctx := context.Background()

	// Создаем тестовый пост
	newPost, err := resolver.Mutation().CreatePost(ctx, "Tестовое название", "Контент", "Автор", true)
	assert.NoError(t, err, "Ошибка при создании поста")
	assert.NotNil(t, newPost, "Созданный пост не должен быть nil")

	// Добавляем коммент к этому посту
	newComment, err := resolver.Mutation().AddComment(ctx, newPost.ID, nil, "Автор коммента", "Тест коммента")
	assert.NoError(t, err, "Ошибка при добавлении комментария")
	assert.NotNil(t, newComment, "Созданный комментарий не должен быть nil")

	// Получаем комменты к этому посту
	comments, err := resolver.Query().Comments(ctx, newPost.ID, 10, 0)
	assert.NoError(t, err, "Ошибка при получении комментариев")
	assert.NotNil(t, comments, "Список комментариев не должен быть nil")
	assert.Equal(t, 1, len(comments), "Должен быть 1 комментарий")
	assert.Equal(t, newComment.ID, comments[0].ID, "ID комментария должен совпадать")
}

// Тест создания поста без комментариев
func TestCreatePostWithoutComments(t *testing.T) {
	resolver := &graph.Resolver{}
	ctx := context.Background()

	newPost, err := resolver.Mutation().CreatePost(ctx, "Тест без комментов", "Контент", "Автор", false)
	assert.NoError(t, err, "Ошибка при создании поста")
	assert.NotNil(t, newPost, "Созданный пост не должен быть nil")
	assert.Equal(t, false, newPost.CommentsAllowed)
}

// Тест добавления комментария к посту, где они запрещены
func TestAddCommentToPostWithDisabledComments(t *testing.T) {
	storage.DB = storage.NewMemoryStorage() // Инициализируем хранилище

	resolver := &graph.Resolver{}
	ctx := context.Background()

	// Создаем пост, в котором запрещены комментарии
	newPost, err := resolver.Mutation().CreatePost(ctx, "Тест без комментов", "Контент", "Автор", false)
	assert.NoError(t, err, "Ошибка при создании поста")
	assert.NotNil(t, newPost, "Созданный пост не должен быть nil")

	// Пытаемся добавить коммента (должна быть ошибка)
	newComment, err := resolver.Mutation().AddComment(ctx, newPost.ID, nil, "Автор", "Коммент")
	assert.Error(t, err, "Ожидалась ошибка при попытке добавить комментарий к посту с отключенными комментариями")
	assert.Nil(t, newComment, "Комментарий не должен быть создан")
	assert.Equal(t, "комментарии к этому посту запрещены", err.Error(), "Текст ошибки должен совпадать")
}

func TestFullProcess(t *testing.T) {
	storage.DB = storage.NewMemoryStorage() // Инициализируем хранилище

	resolver := &graph.Resolver{}
	ctx := context.Background()

	// Создаем пост
	newPost, err := resolver.Mutation().CreatePost(ctx, "Тестовый пост", "Контент", "Автор", true)
	assert.NoError(t, err, "Ошибка при создании поста")
	assert.NotNil(t, newPost, "Созданный пост не должен быть nil")

	// Добавляем комментарий к посту
	newComment, err := resolver.Mutation().AddComment(ctx, newPost.ID, nil, "Тестовый коммент", "Коммент для поста")
	assert.NoError(t, err, "Ошибка при добавлении комментария")
	assert.NotNil(t, newComment, "Созданный комментарий не должен быть nil")

	// Ждем 100 мс (чтобы эмулировать задержку обработки в памяти)
	time.Sleep(100 * time.Millisecond)

	// Получаем пост с комментариями
	retrievedPost, err := resolver.Query().Post(ctx, newPost.ID)
	assert.NoError(t, err, "Ошибка при получении поста")
	assert.NotNil(t, retrievedPost, "Полученный пост не должен быть nil")
	assert.Equal(t, newPost.ID, retrievedPost.ID, "ID поста должен совпадать")

	// Логируем результат для отладки
	t.Logf("Полученные комментарии: %+v", retrievedPost.Comments)

	// Проверяем что комментарий сохранился
	if assert.NotNil(t, retrievedPost.Comments, "Список комментариев не должен быть nil") {
		assert.Equal(t, 1, len(retrievedPost.Comments), "Должен быть 1 комментарий")
		assert.Equal(t, newComment.ID, retrievedPost.Comments[0].ID, "ID комментария должен совпадать")
		assert.Equal(t, "Коммент для поста", retrievedPost.Comments[0].Content, "Текст комментария должен совпадать")
	}
}
func TestPaginationComments(t *testing.T) {
	storage.DB = storage.NewMemoryStorage() // Инициализируем in-memory хранилище

	resolver := &graph.Resolver{}
	ctx := context.Background()

	// Создаем пост
	newPost, err := resolver.Mutation().CreatePost(ctx, "Pagination Test Post", "Test Content", "Test Author", true)
	assert.NoError(t, err, "Ошибка при создании поста")
	assert.NotNil(t, newPost, "Созданный пост не должен быть nil")

	// Добавляем 5 комментов
	for i := 1; i <= 5; i++ {
		_, err := resolver.Mutation().AddComment(ctx, newPost.ID, nil, fmt.Sprintf("Author %d", i), fmt.Sprintf("Comment %d", i))
		assert.NoError(t, err, "Ошибка при добавлении комментария")
	}

	// Ждем чтобы эмулировать задержку сохранения данных
	time.Sleep(100 * time.Millisecond)

	// Получаем первую страницу (2 коммента)
	firstPage, err := resolver.Query().Comments(ctx, newPost.ID, 2, 0)
	assert.NoError(t, err, "Ошибка при получении первой страницы комментариев")
	assert.Equal(t, 2, len(firstPage), "На первой странице должно быть 2 комментария")

	// Получаем вторую страницу (следующие 2 коммента)
	secondPage, err := resolver.Query().Comments(ctx, newPost.ID, 2, 2)
	assert.NoError(t, err, "Ошибка при получении второй страницы комментариев")
	assert.Equal(t, 2, len(secondPage), "На второй странице должно быть 2 комментария")

	// Проверка на то, что данные не совпадают
	assert.NotEqual(t, firstPage[0].ID, secondPage[0].ID, "Первая и вторая страницы не должны содержать одинаковые комментарии")
	assert.NotEqual(t, firstPage[1].ID, secondPage[1].ID, "Первая и вторая страницы не должны содержать одинаковые комментарии")

	// Получаем третью страницу (оставшийся 1 коммент)
	thirdPage, err := resolver.Query().Comments(ctx, newPost.ID, 2, 4)
	assert.NoError(t, err, "Ошибка при получении третьей страницы комментариев")
	assert.Equal(t, 1, len(thirdPage), "На третьей странице должен быть 1 комментарий")

	// Убеждаемся, что это последний комментарий
	assert.NotEqual(t, thirdPage[0].ID, firstPage[0].ID, "Третий комментарий должен быть уникальным")
	assert.NotEqual(t, thirdPage[0].ID, secondPage[0].ID, "Третий комментарий должен быть уникальным")

	// Получаем страницу за пределами доступных комментариев (должно быть пусто)
	emptyPage, err := resolver.Query().Comments(ctx, newPost.ID, 2, 6)
	assert.NoError(t, err, "Ошибка при получении пустой страницы")
	assert.Equal(t, 0, len(emptyPage), "Пустая страница должна возвращать 0 комментариев")
}
func TestNestedComments(t *testing.T) {
	storage.DB = storage.NewMemoryStorage() // Используем in-memory хранилище

	resolver := &graph.Resolver{}
	ctx := context.Background()

	// Создаем пост
	post, err := resolver.Mutation().CreatePost(ctx, "Teстовый пост", "Контент", "Aвтор", true)
	assert.NoError(t, err)
	assert.NotNil(t, post)

	// Добавляем корневой комментарий
	rootComment, err := resolver.Mutation().AddComment(ctx, post.ID, nil, "User1", "Корневой комментарий")
	assert.NoError(t, err)
	assert.NotNil(t, rootComment)

	// Добавляем вложенный комментарий
	childComment1, err := resolver.Mutation().AddComment(ctx, post.ID, &rootComment.ID, "User2", "Вложенный комментарий 1")
	assert.NoError(t, err)
	assert.NotNil(t, childComment1)

	// Добавляем еще один уровень вложенности
	childComment2, err := resolver.Mutation().AddComment(ctx, post.ID, &childComment1.ID, "User3", "Вложенный комментариий 2")
	assert.NoError(t, err)
	assert.NotNil(t, childComment2)

	// Проверяем, что вложенные комментарии есть
	comments, err := resolver.Query().Comments(ctx, post.ID, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(comments), "Должно быть 3 комментария (включая вложенные)")
}
func TestCommentSubscription(t *testing.T) {
	storage.DB = storage.NewMemoryStorage() // In-memory хранилище

	resolver := &graph.Resolver{}
	ctx := context.Background()

	// Создаем пост
	post, err := resolver.Mutation().CreatePost(ctx, "Тест подписки", "Контент", "Автор", true)
	assert.NoError(t, err)
	assert.NotNil(t, post)

	// Подписываемся на комментарии
	commentChan, err := resolver.Subscription().CommentAdded(ctx, post.ID)
	assert.NoError(t, err)

	// Добавляем комментарий
	go func() {
		time.Sleep(1 * time.Second) // Ждем 1 секунду, имитируем задержку
		_, _ = resolver.Mutation().AddComment(ctx, post.ID, nil, "User", "Новый коммент для теста")
	}()

	// Ждем комментарий из подписки
	select {
	case newComment := <-commentChan:
		assert.NotNil(t, newComment, "Комментарий из подписки не должен быть nil")
		assert.Equal(t, "Новый коммент для теста", newComment.Content, "Содержимое комментария должно совпадать")
	case <-time.After(3 * time.Second):
		t.Fatal("Не получили комментарий через подписку")
	}
}
