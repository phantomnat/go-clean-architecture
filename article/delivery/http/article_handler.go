package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/phantomnat/go-clean-architecture/domain"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ResponseError represents the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ArticleHandle represents the http handler for article
type ArticleHandler struct {
	ArticleUsecase domain.ArticleUsecase
}

func NewArticleHttpHandler(e *gin.Engine, au domain.ArticleUsecase) {
	handler := &ArticleHandler{
		ArticleUsecase: au,
	}

	e.GET("/articles", handler.FetchArticle)
	e.GET("/article/:id", handler.GetByID)
}

// FetchArticle will fetch the article based on given params
func (a *ArticleHandler) FetchArticle(c *gin.Context) {
	n := c.Query("num")
	num, _ := strconv.Atoi(n)

	cursor := c.Query("cursor")

	ctx, cancel := context.WithCancel(c)
	defer cancel()

	listAr, nextCursor, err := a.ArticleUsecase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		c.AbortWithStatusJSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	c.Header("X-Cursor", nextCursor)
	c.JSON(http.StatusOK, listAr)
}

// GetByID returns article by given id
func (a *ArticleHandler) GetByID(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, domain.ErrNotFound.Error())
		return
	}

	id := int64(i)
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	ar, err := a.ArticleUsecase.GetByID(ctx, id)
	if err != nil {
		c.AbortWithStatusJSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ar)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case domain.ErrInternalServer:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrAlreadyExist:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
