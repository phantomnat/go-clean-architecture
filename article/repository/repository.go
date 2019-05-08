package repository

import models "github.com/phantomnat/go-clean-architecture/article"

type ArticleRepository interface {
	Fetch(cursor string, num int64) ([]*models.Article, error)
	GetByID(id int64) (*models.Article, error)
	GetByTitle(title string) (*models.Article, error)
	Update(ar *models.Article) (*models.Article, error)
	Store(ar *models.Article) (int64, error)
	Delete(id int64) (bool, error)
}
