package router

import "github.com/gin-gonic/gin"

func Initrouter() {
	router := gin.Default()
	router.GET("/article/id/:id", getArticleByID)
	router.GET("/article/title/:title", getArticleByTitle)
	router.GET("/article", getArticleByCreatedTime)
	router.POST("/article", createArticle)
	router.GET("/health", health)
	router.PUT("/article/:id", updateArticle)
	router.DELETE("/article/:id", deleteArticle)
	router.Run("0.0.0.0:8080")
}
