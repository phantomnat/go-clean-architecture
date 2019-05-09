package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phantomnat/go-clean-architecture/article"
	"github.com/phantomnat/go-clean-architecture/article/usecase"
	"github.com/sirupsen/logrus"
)

type ArticleHandler struct {
	ArticleUsecase usecase.ArticleUsecase
}

func (a *ArticleHandler) FetchArticle(c *gin.Context) {
	n := c.Query("num")
	num, _ := strconv.Atoi(n)

	cursor := c.Query("cursor")

	listAr, nextCursor, err := a.ArticleUsecase.Fetch(cursor, int64(num))
	if err != nil {
		c.AbortWithStatusJSON(getStatusCode(err), err)
		return
	}

	c.Header("X-Cursor", nextCursor)
	c.JSON(getStatusCode(nil), listAr)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case article.ErrInternalServer:
		return http.StatusInternalServerError
	case article.ErrNotFound:
		return http.StatusNotFound
	case article.ErrAlreadyExist:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
