package usecase_test

import (
	"testing"

	"github.com/phantomnat/go-clean-architecture/article"
)

func TestFetch(t *testing.T) {
	//mockArticleRepo := new(mocks.A)
	mockArticle := &article.Article{
		Title:   "hello",
		Content: "content",
	}

	mockListArticle := make([]*article.Article, 0)
	mockListArticle = append(mockListArticle, mockArticle)
}
