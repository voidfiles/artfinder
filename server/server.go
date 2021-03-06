package server

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jmoiron/sqlx"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/voidfiles/artarchive/config"
	"github.com/voidfiles/artarchive/logging"
	"github.com/voidfiles/artarchive/slides"
	"github.com/voidfiles/artarchive/storage"
)

type ContextFunc func(c RequestContext)

func bind(f ContextFunc) func(*gin.Context) {
	return func(c *gin.Context) {
		f(c)
	}
}

func Serve() {
	binding.Validator = new(defaultValidator)
	appConfig := config.NewAppConfig()
	logger := logging.NewLogger(false, os.Stdout)

	db, err := sqlx.Connect(appConfig.Database.Type, appConfig.Database.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	sss := s3.New(sess)

	ginLogger := logging.MustNewGinLogger(logger, "gin")
	router := gin.New()

	router.Use(ginLogger.Logger())
	router.Use(cors.AllowAll())
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": appConfig.AuthPassword,
	}))

	slidesDBStorage := storage.MustNewItemStorage(db)
	slideS3Storage := slides.NewSlideStorage(sss, appConfig.Bucket, appConfig.Version)
	handlers := MustNewServerHandlers(logger, slideS3Storage, slidesDBStorage)

	router.GET("/render/slide/:key", bind(handlers.RenderSlide))
	authorized.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
	authorized.GET("/slides/", bind(handlers.ListSlides))
	authorized.GET("/slides/:key", bind(handlers.GetSlide))
	authorized.PUT("/slides/:key", bind(handlers.UpdateSlide))
	authorized.GET("/sites/", bind(handlers.ListSites))

	router.Run(":" + appConfig.Port)
}
