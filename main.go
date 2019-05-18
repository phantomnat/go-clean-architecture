package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/phantomnat/go-clean-architecture/article/delivery/http"
	articleRepo "github.com/phantomnat/go-clean-architecture/article/repository/mysql"
	"github.com/phantomnat/go-clean-architecture/article/usecase"
	authorRepo "github.com/phantomnat/go-clean-architecture/author/repository/mysql"
	"github.com/phantomnat/go-clean-architecture/config/env"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var config env.Config

func init() {
	config = env.NewViperConfig()

	if config.GetBool("debug") {
		fmt.Println("service run on DEBUG mode")
	}
}

func main() {
	dbHost := config.GetString("database.host")
	dbPort := config.GetString("database.port")
	dbUser := config.GetString("database.user")
	dbName := config.GetString("database.name")
	dbPass := config.GetString("database.pass")
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Bangkok")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil && config.GetBool("debug") {
		logrus.Error(err)
	}
	err = dbConn.Ping()
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
	defer func() {
		err = dbConn.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	router := gin.Default()
	authoreRepo := authorRepo.NewMysqlAuthorRepository(dbConn)
	articleRepo := articleRepo.NewMysqlArticleRepository(dbConn)
	timeoutContext := time.Second * 2
	au := usecase.NewArticleUseCase(articleRepo, authoreRepo, timeoutContext)

	http.NewArticleHttpHandler(router, au)

	router.Run()
}
