package endpoints

import (
	"github.com/gin-gonic/gin"

	"opentunnel/server/config"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	authHandler := NewAuthHandler(cfg)

	router.Static("/static", "./web")
	router.GET("/login", func(c *gin.Context) {
		c.File("./web/login.html")
	})

	api := router.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.GET("/validate", authHandler.ValidateToken)
	}

	return router
}
