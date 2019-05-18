package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/phantomnat/go-clean-architecture/article/usecase"
	"github.com/phantomnat/go-clean-architecture/domain"
	"github.com/phantomnat/go-clean-architecture/domain/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetch(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "hello",
		Content: "content",
	}

	mockListArticle := make([]domain.Article, 0)
	mockListArticle = append(mockListArticle, mockArticle)
	//mockListArticle = append(mockListArticle, domain.Article{
	//	Title:   "hello-2",
	//	Content: "content 2",
	//})

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
			Return(mockListArticle, "next-cursor", nil).Once()
		mockAuthor := domain.Author{
			ID:   1,
			Name: "Iron Man",
		}
		mockAuthorRepo := new(mocks.AuthorRepository)
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := usecase.NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)
		cursorExpected := "next-cursor"

		assert.Equal(t, cursorExpected, nextCursor)
		assert.NotEmpty(t, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListArticle))

		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int64")).
			Return(nil, "", errors.New("unexpected error")).Once()
		mockAuthorRepo := new(mocks.AuthorRepository)
		u := usecase.NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "hello",
		Content: "content",
	}
	mockAuthor := domain.Author{
		ID:   1,
		Name: "Iron Man",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockArticle, nil).Once()
		mockAuthorRepo := new(mocks.AuthorRepository)
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := usecase.NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(domain.Article{}, errors.New("unexpected error")).Once()
		mockAuthorRepo := new(mocks.AuthorRepository)
		u := usecase.NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		a, err := u.GetByID(context.TODO(), mockArticle.ID)
		assert.Error(t, err)
		assert.Equal(t, domain.Article{}, a)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "hello",
		Content: "content",
	}

	t.Run("success", func(t *testing.T) {
		tempMockArticle := mockArticle
		tempMockArticle.ID = 0
		mockArticleRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).
			Return(domain.Article{}, domain.ErrNotFound).Once()
		mockArticleRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.Article")).
			Return(nil).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := usecase.NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)

		err := u.Store(context.TODO(), &tempMockArticle)

		assert.NoError(t, err)
		assert.Equal(t, mockArticle.Title, tempMockArticle.Title)
		mockArticleRepo.AssertExpectations(t)
	})

	t.Run("error existing title", func(t *testing.T) {
		existingArticle := mockArticle
		mockArticleRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).
			Return(existingArticle, nil).Once()
		mockAuthor := domain.Author{
			ID:   1,
			Name: "Iron Man",
		}

		mockAuthorRepo := new(mocks.AuthorRepository)
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(mockAuthor, nil).Once()

		u := usecase.NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		err := u.Store(context.TODO(), &mockArticle)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)

	})
}

func TestDelete(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "hello",
		Content: "content",
	}
	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockArticle, nil).Once()
		mockArticleRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := usecase.NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockArticle.ID)
		assert.NoError(t, err)

		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("article is not exist", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(domain.Article{}, nil).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := usecase.NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
	t.Run("error happens in db", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).
			Return(domain.Article{}, errors.New("unexpected error")).Once()
		mockAuthorRepo := new(mocks.AuthorRepository)
		u := usecase.NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "hello",
		Content: "content",
		ID:      23,
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Update", mock.Anything, &mockArticle).Return(nil).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := usecase.NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)
		err := u.Update(context.TODO(), &mockArticle)

		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)
	})
}
