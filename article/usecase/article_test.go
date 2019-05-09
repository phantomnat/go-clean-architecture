package usecase_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/phantomnat/go-clean-architecture/article/usecase"
	"github.com/stretchr/testify/assert"

	"github.com/phantomnat/go-clean-architecture/author"
	mocksAuthor "github.com/phantomnat/go-clean-architecture/author/repository/mocks"

	"github.com/phantomnat/go-clean-architecture/article/repository/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/phantomnat/go-clean-architecture/article"
)

func TestFetch(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := &article.Article{
		Title:   "hello",
		Content: "content",
	}

	mockListArticle := make([]*article.Article, 0)
	mockListArticle = append(mockListArticle, mockArticle)
	mockListArticle = append(mockListArticle, &article.Article{
		Title:   "hello-2",
		Content: "content 2",
	})

	mockArticleRepo.On("Fetch", mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(mockListArticle, nil)

	mockAuthorRepo := new(mocksAuthor.AuthorRepository)
	mockAuthor := &author.Author{
		ID:   1,
		Name: "Iron Man",
	}
	mockAuthorRepo.On("GetByID", mock.AnythingOfType("int64")).Return(mockAuthor, nil)

	u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo)
	num := int64(2)
	cursor := "12"
	list, nextCursor, err := u.Fetch(cursor, num)
	cursorExpected := strconv.Itoa(int(mockArticle.ID))

	assert.Equal(t, cursorExpected, nextCursor)
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
	assert.Len(t, list, len(mockListArticle))

	mockArticleRepo.AssertCalled(t, "Fetch", mock.AnythingOfType("string"), mock.AnythingOfType("int64"))
}

func TestFetchError(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)

	mockArticleRepo.On("Fetch", mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(nil, errors.New("unexpected error"))

	mockAuthorRepo := new(mocksAuthor.AuthorRepository)

	u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo)
	num := int64(1)
	cursor := "12"
	list, nextCursor, err := u.Fetch(cursor, num)

	assert.Empty(t, nextCursor)
	assert.Error(t, err)
	assert.Len(t, list, 0)

	mockArticleRepo.AssertCalled(t, "Fetch", mock.AnythingOfType("string"), mock.AnythingOfType("int64"))
}

func TestGetByID(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := &article.Article{
		ID:      1,
		Title:   "hello",
		Content: "content",
	}

	mockArticleRepo.On("GetByID", mock.AnythingOfType("int64")).Return(mockArticle, nil)
	defer mockArticleRepo.AssertCalled(t, "GetByID", mock.AnythingOfType("int64"))

	mockAuthorRepo := new(mocksAuthor.AuthorRepository)
	mockAuthor := &author.Author{
		ID:   1,
		Name: "Iron Man",
	}
	mockAuthorRepo.On("GetByID", mock.AnythingOfType("int64")).Return(mockAuthor, nil)

	u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo)
	a, err := u.GetByID(mockArticle.ID)

	assert.NoError(t, err)
	assert.NotNil(t, a)
}

func TestStore(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := article.Article{
		ID:      1,
		Title:   "hello",
		Content: "content",
	}

	// set to 0 because this is test from Client, and ID is an AutoIncrement
	tempMockArticle := mockArticle
	tempMockArticle.ID = 0

	mockArticleRepo.On("GetByTitle", mock.AnythingOfType("string")).Return(nil, article.ErrNotFound)
	mockArticleRepo.On("Store", mock.AnythingOfType("*article.Article")).Return(mockArticle.ID, nil)
	defer mockArticleRepo.AssertCalled(t, "GetByTitle", mock.AnythingOfType("string"))
	defer mockArticleRepo.AssertCalled(t, "Store", mock.AnythingOfType("*article.Article"))

	mockAuthorRepo := new(mocksAuthor.AuthorRepository)
	u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo)

	a, err := u.Store(&tempMockArticle)

	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, mockArticle.Title, tempMockArticle.Title)

}

func TestDelete(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := article.Article{
		ID:      1,
		Title:   "hello",
		Content: "content",
	}
	mockArticleRepo.On("GetByID", mock.AnythingOfType("int64")).Return(&mockArticle, article.ErrNotFound)
	defer mockArticleRepo.AssertCalled(t, "GetByID", mock.AnythingOfType("int64"))

	mockArticleRepo.On("Delete", mock.AnythingOfType("int64")).Return(true, nil)
	defer mockArticleRepo.AssertCalled(t, "Delete", mock.AnythingOfType("int64"))

	mockAuthorRepo := new(mocksAuthor.AuthorRepository)
	u := usecase.NewArticleUsecase(mockArticleRepo, mockAuthorRepo)

	a, err := u.Delete(mockArticle.ID)

	assert.NoError(t, err)
	assert.True(t, a)
}
