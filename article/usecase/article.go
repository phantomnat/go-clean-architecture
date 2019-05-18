package usecase

import (
	"context"
	"time"

	"github.com/phantomnat/go-clean-architecture/domain"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type articleUsecase struct {
	articleRepo    domain.ArticleRepository
	authorRepo     domain.AuthorRepository
	contextTimeout time.Duration
}

var _ domain.ArticleUsecase = &articleUsecase{}

type AuthorChannel struct {
	Author *domain.Author
	Error  error
}

// NewArticleUseCase will create new articleUsecase object representation of domain.ArticleUseCase interface
func NewArticleUseCase(article domain.ArticleRepository, author domain.AuthorRepository, timeout time.Duration) domain.ArticleUsecase {
	return &articleUsecase{
		articleRepo:    article,
		authorRepo:     author,
		contextTimeout: timeout,
	}
}

func (a *articleUsecase) fillAuthorDetails(c context.Context, data []domain.Article) ([]domain.Article, error) {
	g, ctx := errgroup.WithContext(c)

	mapAuthors := map[int64]domain.Author{}

	for _, article := range data {
		mapAuthors[article.Author.ID] = domain.Author{}
	}

	// Using goroutine to fetch the author's detail
	chanAuthor := make(chan domain.Author)
	for authorID := range mapAuthors {
		aID := authorID
		g.Go(func() error {
			res, err := a.authorRepo.GetByID(ctx, aID)
			if err != nil {
				return err
			}
			chanAuthor <- res
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}
		close(chanAuthor)
	}()

	for author := range chanAuthor {
		if author != (domain.Author{}) {
			mapAuthors[author.ID] = author
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// merge the author's data
	for index, item := range data {
		if author, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = author
		}
	}

	return data, nil
}

func (a *articleUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillAuthorDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	return
}

func (a *articleUsecase) GetByID(c context.Context, id int64) (res domain.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Article{}, err
	}

	res.Author = resAuthor
	return res, nil
}

func (a *articleUsecase) Update(c context.Context, ar *domain.Article) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *articleUsecase) GetByTitle(c context.Context, title string) (res domain.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Article{}, err
	}

	res.Author = resAuthor
	return res, nil
}

func (a *articleUsecase) Store(c context.Context, ar *domain.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	existedArticle, _ := a.GetByTitle(ctx, ar.Title)
	if existedArticle != (domain.Article{}) {
		return domain.ErrAlreadyExist
	}

	err = a.articleRepo.Store(ctx, ar)
	return
}

func (a *articleUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedArticle == (domain.Article{}) {
		return domain.ErrNotFound
	}
	return a.articleRepo.Delete(ctx, id)
}
