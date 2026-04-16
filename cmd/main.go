package main

import (
	"log"

	"API/internal/APIs"
	"API/internal/config+conn"
	"API/internal/jwt"
	"API/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config_conn.GetConfig()
	rd := config_conn.RedisConn()
	db := config_conn.DbConn()

	config_conn.RunMigrations(cfg.DatabaseUrl)

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, cfg.HMACKey)
	userHandler := user.NewHandler(userService)

	apiRepo := APIs.NewAPRepository(db)
	apiService := APIs.NewAPService(apiRepo, rd)
	apiHandler := APIs.NewAPIHandler(apiService)

	router := gin.Default()

	userGroup := router.Group("/user")
	{
		userGroup.POST("/login", userHandler.Login)
		userGroup.POST("/register", userHandler.Register)
	}

	apiGroup := router.Group("/api")
	apiGroup.Use(jwt.Authorize(cfg.HMACKey))
	{
		apiGroup.GET("/", apiHandler.GetAll)
		apiGroup.GET("/:id", apiHandler.GetByID)
		apiGroup.POST("/", apiHandler.Create)
		apiGroup.PATCH("/:id", apiHandler.Update)
		apiGroup.DELETE("/:id", apiHandler.Delete)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
