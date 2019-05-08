package usecase

import (
	"strconv"
	"time"

	"github.com/phantomnat/go-clean-architecture/article"
	repoArticle "github.com/phantomnat/go-clean-architecture/article/repository"
	"github.com/phantomnat/go-clean-architecture/author"
	repoAuthor "github.com/phantomnat/go-clean-architecture/author/repository"
	"github.com/sirupsen/logrus"
)

type IArticleUsecase interface {
	Fetch(cursor string, num int64) ([]*article.Article, string, error)
	GetByID(id int64) (*article.Article, error)
	Update(ar *article.Article) (*article.Article, error)
	GetByTitle(title string) (*article.Article, error)
	Store(ar *article.Article) (*article.Article, error)
	Delete(id int64) (bool, error)
}

var _ IArticleUsecase = &ArticleUsecase{}

type ArticleUsecase struct {
	articleRepo repoArticle.ArticleRepository
	authorRepo  repoAuthor.AuthorRepository
}

type AuthorChannel struct {
	Author *author.Author
	Error  error
}

func NewArticleUsecase(article repoArticle.ArticleRepository, author repoAuthor.AuthorRepository) IArticleUsecase {
	return &ArticleUsecase{
		articleRepo: article,
		authorRepo:  author,
	}
}

func (a *ArticleUsecase) getAuthorDetail(item *article.Article, ch chan AuthorChannel) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Debug("Recovered in ", r)
		}
	}()

	res, err := a.authorRepo.GetByID(item.Author.ID)
	holder := AuthorChannel{
		Author: res,
		Error:  err,
	}
	ch <- holder
}

func (a *ArticleUsecase) getAuthorDetails(data []*article.Article) ([]*article.Article, error) {
	ch := make(chan AuthorChannel)
	defer close(ch)
	existingAuthorMap := make(map[int64]bool)
	totalCall := 0
	for _, item := range data {
		if _, ok := existingAuthorMap[item.Author.ID]; !ok {
			existingAuthorMap[item.Author.ID] = true
			go a.getAuthorDetail(item, ch)
		}
		totalCall++
	}

	mapAuthor := make(map[int64]*author.Author)
	totalGorutineCalled := len(existingAuthorMap)
	for i := 0; i < totalGorutineCalled; i++ {
		select {
		case a := <-ch:
			if a.Error == nil && a.Author != nil {
				mapAuthor[a.Author.ID] = a.Author
			}
		case <-time.After(time.Second * 1):
			logrus.Warn("timeout when calling user detail")
		}
	}

	// merge the author
	for index, item := range data {
		if a, ok := mapAuthor[item.Author.ID]; ok {
			data[index].Author = *a
		}
	}
	return data, nil
}

func (a *ArticleUsecase) Fetch(cursor string, num int64) ([]*article.Article, string, error) {
	if num == 0 {
		num = 10
	}

	listArticle, err := a.articleRepo.Fetch(cursor, num)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""

	listArticle, err = a.getAuthorDetails(listArticle)
	if err != nil {
		return nil, "", err
	}

	if size := len(listArticle); size == int(num) {
		lastID := listArticle[num-1].ID
		nextCursor = strconv.Itoa(int(lastID))
	}

	return listArticle, nextCursor, nil
}

func (a *ArticleUsecase) GetByID(id int64) (*article.Article, error) {
	res, err := a.articleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	resAuthor, err := a.authorRepo.GetByID(res.Author.ID)
	if err != nil {
		return nil, err
	}

	res.Author = *resAuthor
	return res, nil
}

func (a *ArticleUsecase) Update(ar *article.Article) (*article.Article, error) {
	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ar)
}

func (a *ArticleUsecase) GetByTitle(title string) (*article.Article, error) {
	res, err := a.articleRepo.GetByTitle(title)
	if err != nil {
		return nil, err
	}

	resAuthor, err := a.authorRepo.GetByID(res.Author.ID)
	if err != nil {
		return nil, err
	}

	res.Author = *resAuthor
	return res, nil
}

func (a *ArticleUsecase) Store(ar *article.Article) (*article.Article, error) {
	existedArticle, _ := a.GetByTitle(ar.Title)
	if existedArticle != nil {
		return nil, article.ErrAlreadyExist
	}

	id, err := a.articleRepo.Store(ar)
	if err != nil {
		return nil, err
	}

	ar.ID = id
	return ar, nil
}

func (a *ArticleUsecase) Delete(id int64) (bool, error) {
	existedArticle, _ := a.articleRepo.GetByID(id)
	if existedArticle == nil {
		return false, article.ErrNotFound
	}
	return a.articleRepo.Delete(id)
}
