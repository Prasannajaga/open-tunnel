package endpoints

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"opentunnel/server/config"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	authHandler := NewAuthHandler(cfg)
	// router.LoadHTMLFiles("../web/login.html")

	router.Static("/static", "web")
	router.GET("/login", func(c *gin.Context) {
		fmt.Print("login coming")
		c.HTML(200, "login.html", gin.H{})
	})

	api := router.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.GET("/validate", authHandler.ValidateToken)

		// CLI Auth Flow
		cliAuth := api.Group("/cli-auth")
		{
			cliAuth.POST("/start", authHandler.StartCLIAuth)
			cliAuth.GET("/poll/:session_id", authHandler.PollCLIAuth)
			cliAuth.POST("/complete/:session_id", authHandler.CompleteCLIAuth)
		}
	}

	return router
}
